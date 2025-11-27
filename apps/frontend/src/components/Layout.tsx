import React from 'react';

const Layout = ({ children }: { children: React.ReactNode }) => {
  return (
    // [수정] h-screen 대신 h-full을 사용하여 부모(#root)의 높이를 꽉 채웁니다.
    // w-full: 너비 100%
    // overflow-hidden: 자식 요소가 넘쳐도 스크롤바 생성 방지
    <main className="h-full w-full bg-background overflow-hidden flex flex-col">
      {children}
    </main>
  );
};

export default Layout;