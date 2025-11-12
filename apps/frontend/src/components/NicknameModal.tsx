import { useState } from 'react';
import useChatStore from '../store/chatStore';

const NicknameModal = () => {
  const [nicknameInput, setNicknameInput] = useState('');
  // login 액션은 이제 사용하지 않으므로, setNickname으로 교체하거나
  // 더 명확한 이름의 액션을 스토어에 추가할 수 있습니다.
  // 여기서는 login을 그대로 사용하되, 토큰 저장 로직을 추가합니다.
  const { login } = useChatStore();
  const [isLoading, setIsLoading] = useState(false); // 로딩 상태 추가

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (nicknameInput.trim().length < 2 || isLoading) return;

    setIsLoading(true);

    try {
      // 브라우저 고유 ID 생성 (localStorage에 없으면 새로 생성)
      let anonymousId = localStorage.getItem('anonymousId');
      if (!anonymousId) {
        anonymousId = crypto.randomUUID();
        localStorage.setItem('anonymousId', anonymousId);
      }

      // 백엔드 API 호출
      const response = await fetch('http://localhost:8080/api/session', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          anonymousId: anonymousId,
          nickname: nicknameInput.trim(),
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to create session');
      }

      const data = await response.json();
      const { token } = data;

      // 받아온 JWT 토큰을 localStorage에 저장
      localStorage.setItem('sessionToken', token);

      // Zustand 스토어의 login 액션에 닉네임 대신 토큰을 전달!
      login(token);

    } catch (error) {
      console.error('Session creation failed:', error);
      alert('세션 생성에 실패했습니다. 잠시 후 다시 시도해주세요.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <h2>채팅에 사용할 닉네임을 입력하세요</h2>
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            value={nicknameInput}
            onChange={(e) => setNicknameInput(e.target.value)}
            placeholder="닉네임 (2~15자)"
            minLength={2}
            maxLength={15}
            required
          />
          <button type="submit">입장하기</button>
        </form>
      </div>
    </div>
  );
};

export default NicknameModal;