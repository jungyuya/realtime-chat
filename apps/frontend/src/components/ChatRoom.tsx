import Header from './Header';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import { useWebSocket } from '../hooks/useWebSocket';

const ChatRoom = () => {
  const { sendMessage } = useWebSocket();

  return (
    // [수정]
    // 기존의 max-w-md, my-8, rounded-2xl, shadow-lg 등을 모두 제거합니다.
    // h-full: 부모(Layout)의 높이를 꽉 채움
    // flex flex-col: 헤더-목록-입력을 수직 배치
    <div className="w-full h-full bg-surface flex flex-col">
      <Header />
      
      {/* 메시지 목록 영역 */}
      <div className="flex-grow overflow-y-auto scrollbar-hide"> 
        <MessageList />
      </div>
      
      <MessageInput onSendMessage={sendMessage} />
    </div>
  );
};

export default ChatRoom;