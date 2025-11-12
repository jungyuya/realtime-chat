import { useEffect, useRef } from 'react';
import useChatStore from '../store/chatStore';

const WEBSOCKET_URL = 'ws://localhost:8080/ws';

// 서버로부터 받는 메시지의 타입을 정의합니다. (hub.go의 Message 구조체와 일치)
interface ReceivedMessage {
  content: string;
  senderId: string;
  senderNickname: string;
}

export const useWebSocket = () => {
  // 스토어에서 필요한 상태와 액션을 가져옵니다.
  const { addMessage, anonymousId, nickname } = useChatStore();
  // ws.current는 WebSocket 인스턴스를 저장하며, 리렌더링을 유발하지 않습니다.
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('sessionToken');
    if (!token) {
      // 토큰이 없으면 아무것도 하지 않습니다.
      return;
    }

    // 1. WebSocket URL에 토큰을 쿼리 파라미터로 추가합니다.
    const socket = new WebSocket(`${WEBSOCKET_URL}?token=${token}`);

    socket.onopen = () => {
      console.log('WebSocket connected');
    };

    socket.onmessage = (event) => {
      // 2. 수신된 JSON 문자열을 파싱합니다.
      const receivedMsg: ReceivedMessage = JSON.parse(event.data);

      // 3. 메시지 중복 제거: 내가 보낸 메시지는 무시합니다.
      if (receivedMsg.senderId === anonymousId) {
        return; // 아무것도 하지 않고 함수 종료
      }

      // 4. 다른 사람이 보낸 메시지만 스토어에 추가합니다.
      const newMessage = {
        id: Date.now(), // 임시 ID
        text: receivedMsg.content,
        sender: receivedMsg.senderNickname, // 실제 닉네임 사용
        isMe: false, // 내가 보낸 것이 아님
      };
      addMessage(newMessage);
    };

    socket.onclose = () => {
      console.log('WebSocket disconnected');
      // TODO: 재연결 로직
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      socket.close();
    };

    ws.current = socket;

    // 컴포넌트가 언마운트될 때 실행되는 정리(cleanup) 함수
    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
    // anonymousId가 변경될 때마다 이 useEffect를 다시 실행하여
    // 새로운 ID로 재연결하도록 의존성 배열에 추가합니다.
  }, [anonymousId, addMessage]);

  // 메시지를 보내는 함수
  const sendMessage = (message: string) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(message);
      // 내가 보낸 메시지를 즉시 스토어에 추가 (낙관적 업데이트)
      addMessage({
        id: Date.now(),
        text: message,
        sender: nickname || 'Me', // 내 닉네임 사용
        isMe: true, // 내가 보낸 메시지임
      });
    } else {
      console.error('WebSocket is not connected.');
    }
  };

  return { sendMessage };
};