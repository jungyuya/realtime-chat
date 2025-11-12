// 서버로부터 받는 WebSocket 메시지의 원시(raw) 형태
export interface RawMessage {
  content: string;
  senderId: string;
  senderNickname: string;
  avatar: string;
  timestamp: string; // JSON으로는 보통 ISO 문자열 형태로 전달됨
}

// 프론트엔드에서 UI 렌더링을 위해 가공된 메시지 형태
export interface Message extends RawMessage {
  id: number; // 리스트 렌더링을 위한 고유 key
  isMe: boolean;
}

// 사용자 프로필 정보
export interface UserProfile {
  anonymousId: string;
  nickname: string;
  avatar: string;
}