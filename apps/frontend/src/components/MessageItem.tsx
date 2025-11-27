import type { Message } from '../types/chat';
import { format } from 'date-fns';
import Avatar from './Avatar';

interface MessageItemProps {
  message: Message;
}

const MessageItem = ({ message }: MessageItemProps) => {
  const formattedTime = message.timestamp ? format(new Date(message.timestamp), 'p') : '';

  return (
    <div className={`flex gap-3 ${message.isMe ? 'flex-row-reverse' : ''}`}>
      
      <div className="w-10 shrink-0"> {/* shrink-0 추가: 아바타 찌그러짐 방지 */}
        {!message.isMe && <Avatar avatarKey={message.avatar} />}
      </div>

      <div className={`flex flex-col max-w-[75%] ${message.isMe ? 'items-end' : 'items-start'}`}>
        {!message.isMe && (
          <span className="text-sm text-gray-600 mb-1 ml-1">{message.senderNickname}</span>
        )}

        <div className={`flex items-end gap-2 ${message.isMe ? 'flex-row-reverse' : ''}`}>
          {/* [수정] text-left 클래스 추가: 텍스트는 항상 왼쪽 정렬 */}
          <p
            className={`px-4 py-2 rounded-2xl break-words text-left ${
              message.isMe 
                ? 'bg-primary text-white rounded-br-none' 
                : 'bg-gray-100 text-text-dark rounded-bl-none'
            }`}
          >
            {message.content}
          </p>
          <span className="text-xs text-gray-400 whitespace-nowrap mb-1">{formattedTime}</span>
        </div>
      </div>
    </div>
  );
};

export default MessageItem;