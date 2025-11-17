import { useEffect } from 'react';
import useChatStore from './store/chatStore';
import NicknameModal from './components/NicknameModal';
import ChatRoom from './components/ChatRoom';
import Layout from './components/Layout'; // Layout 컴포넌트 임포트
import './App.css'; // App.css는 이제 거의 필요 없지만 일단 둡니다.

function App() {
  const { isAuthenticated, login } = useChatStore();

  useEffect(() => {
    const token = localStorage.getItem('sessionToken');
    if (token) {
      login(token);
    }
  }, [login]);

  return (
    // Layout 컴포넌트로 전체를 감쌉니다.
    <Layout>
      {isAuthenticated ? <ChatRoom /> : <NicknameModal />}
    </Layout>
  );
}

export default App;