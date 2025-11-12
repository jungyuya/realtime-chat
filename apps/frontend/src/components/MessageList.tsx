import useChatStore from '../store/chatStore';
import MessageItem from './MessageItem';

const MessageList = () => {
  const { messages } = useChatStore();

  return (
    <div className="message-list">
      {messages.map((msg) => (
        <MessageItem key={msg.id} message={msg} />
      ))}
    </div>
  );
};

export default MessageList;