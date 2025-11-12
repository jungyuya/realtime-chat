// ...
interface Message {
  id: number;
  text: string;
  sender: string;
  isMe: boolean;
}

const MessageItem = ({ message }: { message: Message }) => {
  // isMe 속성을 직접 사용
  return (
    <div className={`message-item ${message.isMe ? 'me' : 'other'}`}>
      {!message.isMe && <span className="sender-nickname">{message.sender}</span>}
      <p>{message.text}</p>
    </div>
  );
};

export default MessageItem;