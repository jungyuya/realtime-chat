import React from 'react';

const Layout = ({ children }: { children: React.ReactNode }) => {
  return (
    // h-screen: 화면(Iframe) 높이 100%
    // w-full: 너비 100%
    // overflow-hidden: 전체 화면 스크롤 방지 (내부 스크롤만 허용)
    // bg-background: 배경색 유지
    <main className="h-screen w-full bg-background overflow-hidden flex flex-col">
      {children}
    </main>
  );
};

export default Layout;