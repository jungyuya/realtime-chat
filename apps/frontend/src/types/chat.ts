// 서버로부터 받는 WebSocket 메시지의 원시(raw) 형태
export interface RawMessage {
  content: string;
  senderId: string;
  senderNickname: string;
  avatar: string;
  timestamp: string;
  id?: number
}

// 프론트엔드에서 UI 렌더링을 위해 가공된 메시지 형태
export interface Message { 
  id: number;
  content: string;
  senderNickname: string;
  avatar: string;
  timestamp: string;
  isMe: boolean;
}

// 사용자 프로필 정보 타입을 null을 허용하도록 명시적으로 정의
export interface UserProfile {
  anonymousId: string | null;
  nickname: string | null;
  avatar: string | null;
}