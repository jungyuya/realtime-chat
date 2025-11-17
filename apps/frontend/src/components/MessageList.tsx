import { useEffect, useRef, Fragment } from 'react';
import useChatStore from '../store/chatStore';
import MessageItem from './MessageItem';
import { isSameDay, format } from 'date-fns';

const MessageList = () => {
  const { messages } = useChatStore();
  const messagesEndRef = useRef<HTMLDivElement | null>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  return (
    // p-4: 내부 여백 
    // flex flex-col: 자식 요소(MessageItem)를 수직으로 쌓음 
    // gap-y-2: 말풍선 간 간격 
    // overflow-x-hidden : 채팅 입력시 애니메이션으로 인한 가로스크롤 나타나는 현상 숨기기
    <div className="message-list p-4 flex flex-col gap-y-2 overflow-x-hidden">
      {/* TypeScript가 'msg'의 타입을 'Message'로 올바르게 자동 추론합니다. */}
      {messages.map((msg, index) => {
        const prevMsg = messages[index - 1];

        const showDateSeparator = !prevMsg || !isSameDay(new Date(prevMsg.timestamp), new Date(msg.timestamp));

        return (
          <Fragment key={msg.id}>
            {showDateSeparator && (
              <div className="relative my-4">
                <hr className="border-t border-gray-200" />
                <span className="absolute left-1/2 -translate-x-1/2 -top-3 bg-chat-bg px-2 text-xs text-gray-500">
                  {format(new Date(msg.timestamp), 'yyyy년 M월 d일')}
                </span>
              </div>
            )}
            <MessageItem message={msg} />
          </Fragment>
        );
      })}
      <div ref={messagesEndRef} />
    </div>
  );
};

export default MessageList;