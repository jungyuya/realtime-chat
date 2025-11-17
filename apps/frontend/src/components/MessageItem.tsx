import type { Message } from '../types/chat';
import { format } from 'date-fns';
import Avatar from './Avatar';

interface MessageItemProps {
  message: Message;
}

const MessageItem = ({ message }: MessageItemProps) => {
  // 타임스탬프 포맷팅
  const formattedTime = message.timestamp ? format(new Date(message.timestamp), 'p') : '';

  return (
    // [Tailwind 적용]
    // flex gap-3: 아바타와 내용 사이의 가로 간격
    // isMe 값에 따라 flex-row-reverse 클래스를 추가하여 레이아웃을 뒤집음
    <div
      className={`flex gap-3 ${message.isMe
        ? 'flex-row-reverse animate-slide-in-right'
        : 'animate-slide-in-left'
        }`}
    >
      {/* 아바타 (내가 보낸 메시지에는 표시하지 않음) */}
      <div className="w-10"> {/* 아바타의 공간을 항상 차지하도록 빈 div 추가 */}
        {!message.isMe && <Avatar avatarKey={message.avatar} />}
      </div>

      {/* 닉네임과 말풍선, 타임스탬프를 감싸는 컨테이너 */}
      <div className={`flex flex-col max-w-[75%] ${message.isMe ? 'items-end' : 'items-start'}`}>
        {/* 닉네임 (내가 보낸 메시지가 아닐 때만 표시) */}
        {!message.isMe && (
          <span className="text-sm text-gray-600 mb-1">{message.senderNickname}</span>
        )}

        {/* 말풍선과 타임스탬프 */}
        <div className={`flex items-end gap-2 ${message.isMe ? 'flex-row-reverse' : ''}`}>
          {/* 말풍선 */}
          <p
            className={`px-4 py-2 rounded-2xl break-words ${message.isMe
              ? 'bg-primary text-white rounded-br-none'
              : 'bg-gray-100 text-text-dark rounded-bl-none'
              }`}
          >
            {message.content}
          </p>
          {/* 타임스탬프 */}
          <span className="text-xs text-gray-400 whitespace-nowrap">{formattedTime}</span>
        </div>
      </div>
    </div>
  );
};

export default MessageItem;