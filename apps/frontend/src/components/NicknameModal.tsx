import { useState } from 'react';
import useChatStore from '../store/chatStore';


const NicknameModal = () => {
  const { login } = useChatStore();
  const [isLoading, setIsLoading] = useState(false);

  // 이제 버튼 클릭 이벤트 핸들러가 됩니다.
  const handleEnter = async () => {
    if (isLoading) return;
    setIsLoading(true);

    try {
      let anonymousId = localStorage.getItem('anonymousId');
      if (!anonymousId) {
        anonymousId = crypto.randomUUID();
        localStorage.setItem('anonymousId', anonymousId);
      }

      const response = await fetch('http://localhost:8080/api/session', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        // 이제 body에는 anonymousId만 보냅니다.
        body: JSON.stringify({ anonymousId }),
      });

      if (!response.ok) {
        throw new Error('Failed to create session');
      }

      const { token } = await response.json();
      localStorage.setItem('sessionToken', token);
      
      // 스토어의 login 액션에 토큰을 전달합니다.
      // 이 부분은 이전 단계에서 이미 수정되었습니다.
      login(token);

    } catch (error) {
      console.error('Session creation failed:', error);
      alert('세션 생성에 실패했습니다.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <h2>실시간 채팅방에 오신 것을 환영합니다!</h2>
        <p>버튼을 눌러 대화에 참여하세요.</p>
        <button onClick={handleEnter} disabled={isLoading}>
          {isLoading ? '입장 중...' : '입장하기'}
        </button>
      </div>
    </div>
  );
};

export default NicknameModal;