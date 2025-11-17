interface AvatarProps {
  avatarKey: string | null;
}

const Avatar = ({ avatarKey }: AvatarProps) => {
  // 아바타 키가 없거나 null일 경우를 대비해 기본 이미지 경로 설정
  const imageUrl = avatarKey ? `/avatars/${avatarKey}.png` : '/avatars/default.png';

  return (
    <img
      src={imageUrl}
      alt="User avatar"
      className="w-10 h-10 rounded-full object-cover bg-gray-200 border border-gray-300"
    />
  );
};

export default Avatar;