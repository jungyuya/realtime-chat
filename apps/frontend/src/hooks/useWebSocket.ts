import { useEffect, useRef, useCallback } from 'react';
import useChatStore from '../store/chatStore';
import type { RawMessage } from '../types/chat';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export const useWebSocket = () => {
  const { addMessage, anonymousId } = useChatStore();
  const ws = useRef<WebSocket | null>(null);
  const reconnectTimeout = useRef<ReturnType<typeof setTimeout> | null>(null);

  const connect = useCallback(() => {
    if (!anonymousId) return;
    const token = localStorage.getItem('sessionToken');
    if (!token) return;

    // 이미 연결된 상태면 패스
    if (ws.current && (ws.current.readyState === WebSocket.OPEN || ws.current.readyState === WebSocket.CONNECTING)) {
      return;
    }

    const wsUrl = `${API_BASE_URL.replace(/^http/, 'ws')}/ws?token=${token}`;
    console.log("Connecting to WebSocket...");
    
    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log('WebSocket connected');
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
        reconnectTimeout.current = null;
      }
    };

    socket.onmessage = (event) => {
      try {
        const receivedMsg: RawMessage = JSON.parse(event.data);
        // [중요] 스토어의 중복 방지 로직을 믿고 무조건 보냅니다.
        addMessage(receivedMsg, anonymousId);
      } catch (error) {
        console.error("Failed to parse:", error);
      }
    };

    socket.onclose = () => {
      console.log('WebSocket disconnected.');
      ws.current = null;
      // 3초 뒤 재연결 시도
      reconnectTimeout.current = setTimeout(() => connect(), 3000);
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      socket.close();
    };

    ws.current = socket;
  }, [anonymousId, addMessage]);

  // [추가] 화면 가시성(Visibility) 감지 로직
  useEffect(() => {
    const handleVisibilityChange = () => {
      // 사용자가 탭으로 돌아왔고(visible), 연결이 끊겨있다면 즉시 재연결
      if (document.visibilityState === 'visible') {
        if (!ws.current || ws.current.readyState === WebSocket.CLOSED) {
          console.log("Tab is visible, reconnecting...");
          connect();
        }
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    
    // 초기 연결
    connect();

    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      if (reconnectTimeout.current) clearTimeout(reconnectTimeout.current);
      if (ws.current) ws.current.close();
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