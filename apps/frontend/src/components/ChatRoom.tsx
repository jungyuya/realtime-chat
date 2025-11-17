import Header from './Header';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import { useWebSocket } from '../hooks/useWebSocket';

const ChatRoom = () => {
  // useWebSocket 훅에서 sendMessage 함수를 가져옵니다.
  const { sendMessage } = useWebSocket();

  return (
    <div className="w-full max-w-md h-full sm:h-auto sm:my-8 bg-surface rounded-none sm:rounded-2xl shadow-lg flex flex-col overflow-hidden">      <Header />
      {/* MessageList는 남는 공간을 모두 차지해야 합니다. */}
      <div className="flex-grow overflow-y-auto">
        <MessageList />
      </div>
      <MessageInput onSendMessage={sendMessage} />
    </div>
  );
};

export default ChatRoom;