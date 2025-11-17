import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  // [추가] 프로덕션 빌드 옵션을 설정합니다.
  build: {
    sourcemap: false, // 소스맵을 생성하지 않도록 설정
  },
})