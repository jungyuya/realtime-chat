import { create } from 'zustand';
import { jwtDecode } from 'jwt-decode'; // jwt-decode 임포트

// JWT 페이로드의 타입을 정의합니다.
interface JwtPayload {
  anonymousId: string;
  nickname: string;
  exp: number;
  iat: number;
}

// 메시지의 타입을 정의합니다.
interface Message {
  id: number;
  text: string;
  sender: string; // 이제 닉네임이 저장됩니다.
  isMe: boolean; // 내가 보낸 메시지인지 여부
}

// 스토어의 상태(state) 타입을 정의합니다.
interface ChatState {
  messages: Message[];
  isAuthenticated: boolean;
  nickname: string | null;
  anonymousId: string | null; // anonymousId 상태 추가
}

interface ChatActions {
  addMessage: (message: Message) => void;
  login: (token: string) => void; // 이제 토큰을 직접 받습니다.
  logout: () => void;
}

const useChatStore = create<ChatState & ChatActions>((set) => ({
  // 초기 상태는 동일
  messages: [],
  isAuthenticated: false,
  nickname: null,
  anonymousId: null,

  addMessage: (message) =>
    set((state) => ({
      messages: [...state.messages, message],
    })),

  // login 액션의 새로운 구현
  login: (token) => {
    try {
      const decoded = jwtDecode<JwtPayload>(token);
      set({
        isAuthenticated: true,
        nickname: decoded.nickname,
        anonymousId: decoded.anonymousId,
      });
    } catch (error) {
      console.error("Failed to decode token:", error);
    }
  },

  logout: () => {
    localStorage.removeItem('sessionToken'); // 로그아웃 시 토큰도 삭제
    localStorage.removeItem('anonymousId');
    set({ isAuthenticated: false, nickname: null, anonymousId: null, messages: [] });
  },
}));

export default useChatStore;