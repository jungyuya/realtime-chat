import { useEffect, useRef } from 'react';
import useChatStore from '../store/chatStore';
import type { RawMessage } from '../types/chat';

// [수정] 환경 변수로부터 WebSocket URL을 동적으로 생성
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL; 


export const useWebSocket = () => {
  const { addMessage, anonymousId } = useChatStore();
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!anonymousId) {
      return;
    }

    const token = localStorage.getItem('sessionToken');
    if (!token) {
      console.error("Session token not found. Cannot connect to WebSocket.");
      return;
    }

    // [수정] http를 ws로 바꾸고, /ws 경로를 명시적으로 추가합니다.
    const wsUrl = `${API_BASE_URL.replace(/^http/, 'ws')}/ws?token=${token}`;
    const socket = new WebSocket(wsUrl);

    // WebSocket 연결이 성공적으로 열렸을 때 호출되는 이벤트 핸들러입니다.
    socket.onopen = () => {
      console.log('WebSocket connected');
    };

    // 서버로부터 메시지를 수신했을 때 호출되는 이벤트 핸들러입니다.
    // 모든 메시지(내가 보낸 것 포함)는 이 핸들러를 통해 처리됩니다.
    socket.onmessage = (event) => {
      // --- 디버깅 로그 강화 ---
      console.log("Received raw data string from server:", event.data);
      // --------------------
      try {
        const receivedMsg: RawMessage = JSON.parse(event.data);

        // --- 디버깅 로그 강화 ---
        console.log("Parsed message object:", JSON.stringify(receivedMsg, null, 2));
        // --------------------

        if (anonymousId) {
          addMessage(receivedMsg, anonymousId);
        }
      } catch (error) {
        console.error("Failed to parse incoming message:", error);
      }
    };

    // WebSocket 연결이 닫혔을 때 호출되는 이벤트 핸들러입니다.
    socket.onclose = () => {
      console.log('WebSocket disconnected');
      // TODO: 여기에 Exponential Backoff를 이용한 자동 재연결 로직을 구현할 수 있습니다.
    };

    // WebSocket 통신 중 에러가 발생했을 때 호출되는 이벤트 핸들러입니다.
    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      socket.close(); // 에러 발생 시 안전하게 연결을 종료합니다.
    };

    // 생성된 WebSocket 인스턴스를 ws.current에 저장하여 다른 함수(sendMessage)에서 참조할 수 있도록 합니다.
    ws.current = socket;

    // useEffect의 '정리(cleanup)' 함수입니다.
    // 컴포넌트가 언마운트되거나, 의존성 배열의 값이 바뀌어 useEffect가 다시 실행되기 직전에 호출됩니다.
    // 이를 통해 불필요한 WebSocket 연결이 중복으로 생성되는 것을 방지합니다.
    return () => {
      if (ws.current) {
        console.log('Closing WebSocket connection...');
        ws.current.close();
      }
    };
    // 의존성 배열에 anonymousId와 addMessage를 전달합니다.
    // anonymousId가 바뀌면(로그인/로그아웃) 기존 연결을 끊고 새로운 연결을 생성합니다.
  }, [anonymousId, addMessage]);

  // 메시지를 보내는 함수를 정의합니다. 이 함수는 UI 컴포넌트에서 호출됩니다.
  const sendMessage = (messageContent: string) => {
    // WebSocket 연결이 존재하고, 연결 상태(readyState)가 OPEN인지 확인합니다.
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      // [수정] 이제 이 함수는 순수하게 서버로 메시지를 보내는 역할만 합니다.
      // 화면에 메시지를 추가하는 로직(addMessage)은 여기서 호출하지 않습니다.
      ws.current.send(messageContent);
    } else {
      console.error('WebSocket is not connected.');
    }
  };

  // 훅이 외부로 노출하는, 컴포넌트가 사용할 함수를 반환합니다.
  return { sendMessage };
};
