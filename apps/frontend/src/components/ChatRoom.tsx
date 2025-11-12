import { useWebSocket } from '../hooks/useWebSocket';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import useChatStore from '../store/chatStore';

const ChatRoom = () => {
  const { nickname } = useChatStore();
  // sendMessage 함수만 가져옵니다.
  const { sendMessage } = useWebSocket();

  return (
    <div className="chat-room">
      <header>
        <h1>Chat Room</h1>
        <p>Welcome, {nickname}!</p>
        {/* Status 표시는 잠시 제거 */}
      </header>
      <MessageList />
      <MessageInput onSendMessage={sendMessage} />
    </div>
  );
};

export default ChatRoom;