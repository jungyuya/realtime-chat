import { useState, useRef } from 'react';
import type { ChangeEvent, KeyboardEvent } from 'react'; 
import { IoSend } from 'react-icons/io5';

interface MessageInputProps {
  onSendMessage: (message: string) => void;
}

const MessageInput = ({ onSendMessage }: MessageInputProps) => {
  const [input, setInput] = useState('');
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // textarea 높이 자동 조절 함수
  const handleInput = (e: ChangeEvent<HTMLTextAreaElement>) => {
    setInput(e.target.value);
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = 'auto'; // 높이를 초기화
      const scrollHeight = textarea.scrollHeight;
      // 최대 높이를 120px로 제한
      textarea.style.height = `${Math.min(scrollHeight, 120)}px`;
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (input.trim()) {
      onSendMessage(input.trim());
      setInput('');
      // 전송 후 textarea 높이 초기화
      const textarea = textareaRef.current;
      if (textarea) {
        textarea.style.height = 'auto';
      }
    }
  };

  // [추가] 키보드 이벤트 핸들러
  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    // e.key가 'Enter'이고, e.shiftKey가 눌리지 않았을 때 (즉, 그냥 Enter만 눌렀을 때)
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault(); // Enter 키의 기본 동작(줄바꿈)을 막음
      handleSubmit(e as any); // handleSubmit 함수를 호출하여 메시지 전송
    }
    // Shift + Enter를 누르면 기본 동작(줄바꿈)이 그대로 실행됨
  };

  return (
    // [Tailwind 적용]
    <form className="p-4 border-t flex items-center gap-4" onSubmit={handleSubmit}>
      <textarea
        ref={textareaRef}
        value={input}
        onChange={handleInput}
        onKeyDown={handleKeyDown} // [추가] onKeyDown 이벤트 리스너 연결
        placeholder="메시지를 입력하세요..."
        className="flex-grow resize-none overflow-y-auto bg-gray-100 rounded-2xl py-2 px-4 outline-none focus:ring-2 focus:ring-primary"
        rows={1} // 초기 높이는 1줄
      />
      <button
        type="submit"
        // 입력값이 있을 때만 활성화되는 동적 스타일
        className={`p-3 rounded-full text-white transition-colors ${input.trim() ? 'bg-primary hover:bg-primary-hover' : 'bg-gray-300 cursor-not-allowed'
          }`}
        disabled={!input.trim()}
        aria-label="메시지 전송"
      >
        <IoSend />
      </button>
    </form>
  );
};

export default MessageInput;