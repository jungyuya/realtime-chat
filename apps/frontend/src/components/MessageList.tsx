import { useEffect, useRef } from 'react'; // useRef와 useEffect 임포트
import useChatStore from '../store/chatStore';
import MessageItem from './MessageItem';

const MessageList = () => {
  const { messages } = useChatStore();
  // div 요소를 참조할 ref 객체를 생성합니다.
  const messagesEndRef = useRef<HTMLDivElement | null>(null);

  const scrollToBottom = () => {
    // ref가 가리키는 요소(마지막 빈 div)를 화면에 보이도록 스크롤합니다.
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  // messages 배열이 변경될 때마다 scrollToBottom 함수를 호출합니다.
  useEffect(() => {
    scrollToBottom();
  }, [messages]); // 의존성 배열에 messages를 넣어, messages가 바뀔 때만 실행되도록 함

  return (
    <div className="message-list">
      {messages.map((msg) => (
        <MessageItem key={msg.id} message={msg} />
      ))}
      {/* 항상 메시지 목록의 가장 마지막에 위치하는 보이지 않는 요소 */}
      <div ref={messagesEndRef} />
    </div>
  );
};

export default MessageList;