import { useEffect, useRef, useCallback } from 'react';
import useChatStore from '../store/chatStore';
import type { RawMessage } from '../types/chat';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export const useWebSocket = () => {
  // logout 액션을 가져옵니다.
  const { addMessage, anonymousId, logout } = useChatStore();
  const ws = useRef<WebSocket | null>(null);
  const reconnectTimeout = useRef<ReturnType<typeof setTimeout> | null>(null);

  const connect = useCallback(() => {
    // 1. 로그인 상태가 아니면(anonymousId 없음) 연결하지 않음
    if (!anonymousId) return;

    // 2. 토큰 확인
    const token = localStorage.getItem('sessionToken');
    if (!token) {
      console.error("Session token not found. Logging out.");
      // [수정] 토큰이 없으면 스토어 상태도 로그아웃으로 초기화하여 UI를 로그인 화면으로 돌립니다.
      logout();
      return;
    }

    // 3. 이미 연결된 상태면 중복 연결 방지
    if (ws.current && (ws.current.readyState === WebSocket.OPEN || ws.current.readyState === WebSocket.CONNECTING)) {
      return;
    }

    // 4. URL 생성 (http -> ws 변환)
    const wsUrl = `${API_BASE_URL.replace(/^http/, 'ws')}/ws?token=${token}`;
    console.log("Connecting to WebSocket...");
    
    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log('WebSocket connected');
      // 연결 성공 시 재연결 타이머 해제
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
        reconnectTimeout.current = null;
      }
    };

    socket.onmessage = (event) => {
      try {
        const receivedMsg: RawMessage = JSON.parse(event.data);
        // 스토어의 addMessage 액션 호출 (중복 방지 및 isMe 판별은 스토어 내부에서 처리)
        addMessage(receivedMsg, anonymousId);
      } catch (error) {
        console.error("Failed to parse:", error);
      }
    };

    socket.onclose = (event) => {
      console.log('WebSocket disconnected.', event.code, event.reason);
      ws.current = null;

      // [재연결 로직]
      // 연결이 끊어지면 3초 후에 다시 연결을 시도합니다.
      // 만약 토큰이 만료되었다면, 다음 connect() 실행 시 
      // chatStore의 검증 로직이나 위쪽의 !token 체크에 걸려 로그아웃 처리될 것입니다.
      reconnectTimeout.current = setTimeout(() => {
        connect();
      }, 3000);
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      socket.close(); // 에러 발생 시 소켓을 닫아 onclose를 트리거합니다.
    };

    ws.current = socket;
  }, [anonymousId, addMessage, logout]); // 의존성 배열에 logout 추가

  // [추가] 화면 가시성(Visibility) 감지 로직
  // 사용자가 다른 탭에 갔다가 돌아왔을 때 연결이 끊겨있다면 즉시 재연결합니다.
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible') {
        if (!ws.current || ws.current.readyState === WebSocket.CLOSED) {
          console.log("Tab is visible, reconnecting...");
          connect();
        }
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    
    // 컴포넌트 마운트 시 초기 연결 시도
    connect();

    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      if (reconnectTimeout.current) clearTimeout(reconnectTimeout.current);
      if (ws.current) {
        console.log('Closing WebSocket connection...');
        ws.current.close();
      }
    };
  }, [connect]);

  const sendMessage = (messageContent: string) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(messageContent);
    } else {
      console.error('WebSocket not connected. Reconnecting...');
      // 전송 실패 시 즉시 재연결 시도
      connect();
    }
  };

  return { sendMessage };
};