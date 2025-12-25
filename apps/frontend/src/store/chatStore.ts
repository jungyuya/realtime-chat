import { create } from 'zustand';
import { jwtDecode } from 'jwt-decode';
// [해결 1] 'type' 키워드를 사용하여 타입 전용 임포트를 수행합니다.
import type { Message, UserProfile, RawMessage } from '../types/chat';

// JWT 페이로드 타입 정의 (avatar 추가)
interface JwtPayload extends Omit<UserProfile, 'anonymousId'> { // Omit을 사용하여 anonymousId 제외
  anonymousId: string; // anonymousId는 항상 string
  exp: number;
  iat: number;
}

// [해결 2] ChatState가 UserProfile을 직접 포함하도록 하고, 타입 불일치 해결
interface ChatState extends UserProfile {
  messages: Message[];
  isAuthenticated: boolean;
}

interface ChatActions {
  addMessage: (message: RawMessage, currentUserAnonymousId: string) => void;
  login: (token: string) => void;
  logout: () => void;
}

const useChatStore = create<ChatState & ChatActions>((set) => ({
  // 초기 상태: 이제 타입 정의와 완벽하게 일치합니다.
  messages: [],
  isAuthenticated: false,
  nickname: null,
  anonymousId: null,
  avatar: null,

  // addMessage 액션 수정: RawMessage를 받아 가공하도록 변경
  addMessage: (rawMessage, currentUserAnonymousId) =>
    set((state) => {
      // [수정] 중복 메시지 방지 로직 (ID 기반)
      // 이미 스토어에 같은 ID를 가진 메시지가 있다면 추가하지 않고 무시합니다.
      // (서버가 재연결 시 최근 메시지를 다시 보내줄 때 중복을 막기 위함)
      // rawMessage에 id가 없다면(실시간 메시지) Date.now()로 임시 ID를 쓰지만,
      // DB에서 온 메시지는 고유 ID가 있으므로 그것을 기준으로 합니다.

      // DB 메시지인 경우(id가 있음) 중복 체크
      if (rawMessage.id && state.messages.some((m) => m.id === rawMessage.id)) {
        return state;
      }

      const newMessage: Message = {
        id: rawMessage.id || Date.now(), // DB ID가 있으면 쓰고, 없으면 임시 ID
        content: rawMessage.content,
        senderNickname: rawMessage.senderNickname,
        avatar: rawMessage.avatar,
        timestamp: rawMessage.timestamp,
        isMe: rawMessage.senderId === currentUserAnonymousId,
      };

      return { messages: [...state.messages, newMessage] };
    }),

  login: (token) => {
    try {
      const decoded = jwtDecode<JwtPayload>(token);
      
      // [추가] 토큰 만료 시간(exp) 체크
      // exp는 초 단위, Date.now()는 밀리초 단위이므로 변환 필요
      const currentTime = Date.now() / 1000;
      
      if (decoded.exp < currentTime) {
        console.warn("Token expired. Logging out.");
        localStorage.removeItem('sessionToken');
        localStorage.removeItem('anonymousId');
        // 상태를 업데이트하지 않고 리턴하여 로그인 실패 처리
        return; 
      }

      set({
        isAuthenticated: true,
        nickname: decoded.nickname,
        anonymousId: decoded.anonymousId,
        avatar: decoded.avatar,
      });
    } catch (error) {
      console.error("Failed to decode token:", error);
      // 토큰 형식이 잘못된 경우도 삭제
      localStorage.removeItem('sessionToken');
      localStorage.removeItem('anonymousId');
    }
  },

  logout: () => {
    localStorage.removeItem('sessionToken');
    localStorage.removeItem('anonymousId');
    set({
      isAuthenticated: false,
      nickname: null,
      anonymousId: null,
      avatar: null,
      messages: [],
    });
  },
}));

export default useChatStore;