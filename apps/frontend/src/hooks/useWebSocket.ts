import { useEffect, useRef, useCallback } from 'react'; // useCallback 추가
import useChatStore from '../store/chatStore';
import type { RawMessage } from '../types/chat';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export const useWebSocket = () => {
  const { addMessage, anonymousId } = useChatStore();
  const ws = useRef<WebSocket | null>(null);
  // 재연결 타이머를 저장할 ref
  const reconnectTimeout = useRef<NodeJS.Timeout | null>(null);

  // [수정] connect 함수를 useEffect 밖으로 빼고 useCallback으로 감쌉니다.
  // 이렇게 하면 재귀적으로 호출하거나 외부에서 호출하기 용이합니다.
  const connect = useCallback(() => {
    if (!anonymousId) return;

    const token = localStorage.getItem('sessionToken');
    if (!token) {
      console.error("Session token not found.");
      return;
    }

    // 이미 연결되어 있거나 연결 중이면 중복 연결 방지
    if (ws.current && (ws.current.readyState === WebSocket.OPEN || ws.current.readyState === WebSocket.CONNECTING)) {
      return;
    }

    const wsUrl = `${API_BASE_URL.replace(/^http/, 'ws')}/ws?token=${token}`;
    console.log("Attempting to connect WebSocket...");

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
        if (anonymousId) {
          addMessage(receivedMsg, anonymousId);
        }
      } catch (error) {
        console.error("Failed to parse incoming message:", error);
      }
    };

    socket.onclose = () => {
      console.log('WebSocket disconnected. Trying to reconnect...');
      ws.current = null;

      // [핵심] 연결이 끊어지면 3초 후에 다시 연결을 시도합니다. (Exponential Backoff 대신 단순 3초 딜레이 사용)
      reconnectTimeout.current = setTimeout(() => {
        connect();
      }, 3000);
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      socket.close(); // 에러 발생 시 닫으면 onclose가 트리거되어 재연결 시도함
    };

    ws.current = socket;
  }, [anonymousId, addMessage]);

  useEffect(() => {
    // 컴포넌트 마운트 시 연결 시도
    connect();

    // 언마운트 시 정리
    return () => {
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
      }
      if (ws.current) {
        console.log('Closing WebSocket connection...');
        // 언마운트 시에는 재연결을 막기 위해 onclose 핸들러를 비웁니다.
        ws.current.onclose = null;
        ws.current.close();
      }
    };
  }, [connect]);

  const sendMessage = (messageContent: string) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(messageContent);
    } else {
      console.error('WebSocket is not connected. Message not sent.');
      // 연결이 끊겨있을 때 전송을 시도하면 즉시 재연결을 시도할 수도 있습니다.
      connect();
    }
  };

  return { sendMessage };
};