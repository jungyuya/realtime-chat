import useChatStore from '../store/chatStore';
import { IoSearchOutline, IoEllipsisVertical } from "react-icons/io5";
import Avatar from './Avatar'; // [추가] Avatar 컴포넌트를 임포트합니다.

const Header = () => {
  // [수정] 스토어에서 nickname과 함께 avatar 상태를 가져옵니다.
  const { nickname, avatar } = useChatStore();

  return (
    //  shrink-0 클래스를 추가하여, flex 컨테이너 안에서 헤더의 크기가 줄어들지 않도록 합니다.
    <header className="p-4 border-b flex items-center justify-between shrink-0">
      {/* 아바타와 텍스트 정보를 묶는 컨테이너를 추가하고 스타일을 적용합니다. */}
      <div className="flex items-center gap-3">
        {/* Avatar 컴포넌트를 사용하여 현재 사용자의 아바타를 표시합니다. */}
        <Avatar avatarKey={avatar} />
        <div>
          <h1 className="font-bold text-lg text-text-dark leading-tight">{nickname || 'Chat'}</h1>
          <p className="text-xs text-green-500">Online</p>
        </div>
      </div>
      <div className="flex items-center gap-4">
        <button className="text-2xl text-gray-500 hover:text-text-dark">
          <IoSearchOutline />
        </button>
        <button className="text-2xl text-gray-500 hover:text-text-dark">
          <IoEllipsisVertical />
        </button>
      </div>
    </header>
  );
};

export default Header;