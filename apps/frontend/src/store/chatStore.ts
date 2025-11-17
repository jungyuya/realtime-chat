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
      const newMessage: Message = {
        id: Date.now(),
        content: rawMessage.content,
        senderNickname: rawMessage.senderNickname,
        avatar: rawMessage.avatar,
        timestamp: rawMessage.timestamp,
        isMe: rawMessage.senderId === currentUserAnonymousId,
      };

      // =================================================================
      // [검증 2] 빈 말풍선 문제의 원인 찾기 (Zustand Store)
      // =================================================================
      console.group("Zustand addMessage Action");
      console.log("Raw message received by action:", rawMessage);
      console.log("Current user ID for comparison:", currentUserAnonymousId);
      console.log("Is this message from me?:", newMessage.isMe);
      console.log("Final message object to be added to state:", newMessage);
      console.groupEnd();
      // =================================================================

      // [중복 방지 로직]
      // 만약 이 메시지가 내가 보낸 것이고 (isMe: true),
      // 낙관적 업데이트로 인해 이미 스토어에 비슷한 메시지가 있다면 추가하지 않음
      // (간단한 중복 방지: 내용과 isMe 여부가 같은 최근 메시지가 있으면 무시)
      const lastMessage = state.messages[state.messages.length - 1];
      if (
        newMessage.isMe &&
        lastMessage &&
        lastMessage.isMe &&
        lastMessage.content === newMessage.content
      ) {
        // 이미 낙관적 업데이트로 추가된 메시지로 보이므로, 서버로부터 온 것은 무시
        console.log("Duplicate message detected, ignoring server echo.");
        return state; // 상태 변경 없음
      }

      return { messages: [...state.messages, newMessage] };
    }),

  login: (token) => {
    try {
      const decoded = jwtDecode<JwtPayload>(token);
      set({
        isAuthenticated: true,
        nickname: decoded.nickname,
        anonymousId: decoded.anonymousId,
        avatar: decoded.avatar,
      });
    } catch (error) {
      console.error("Failed to decode token:", error);
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