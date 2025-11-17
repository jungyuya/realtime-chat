import React from 'react';

const Layout = ({ children }: { children: React.ReactNode }) => {
  return (
    // [수정] items-center 클래스를 제거합니다.
    // 이제 자식 요소는 컨테이너의 상단부터 채워지게 됩니다.
    <main className="h-screen bg-background flex justify-center">
      {children}
    </main>
  );
};

export default Layout;