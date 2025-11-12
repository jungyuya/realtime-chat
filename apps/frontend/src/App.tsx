import { useEffect } from 'react'; // useEffect 임포트
import useChatStore from './store/chatStore';
import NicknameModal from './components/NicknameModal';
import ChatRoom from './components/ChatRoom';
import './App.css';

function App() {
  const { isAuthenticated, login } = useChatStore();

  useEffect(() => {
    const token = localStorage.getItem('sessionToken');
    if (token) {
      // localStorage에 토큰이 있으면, login 액션에 토큰을 전달하여 상태를 복원
      login(token);
    }
  }, [login]); // login 함수는 한번만 생성되므로 의존성 배열에 추가해도 안전합니다.

  return (
    <div className="App">
      {isAuthenticated ? <ChatRoom /> : <NicknameModal />}
    </div>
  );
}

export default App;