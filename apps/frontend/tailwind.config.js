/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'background': '#f8fbfd',
        'surface': '#ffffff',
        'chat-bg': '#e8f4f9',

        'primary': '#a5d8f0',
        'primary-hover': '#8ecbe7',

        'text-dark': '#405d6b',
        'text-light': '#ffffff',
        'text-assistant': '#5ab0d3',
        'text-notice': '#7a8e98',

        'user-bubble': '#dff1f8',
        'assistant-bubble': '#edf6f9',
      },

      keyframes: {
        // 왼쪽에서 들어오는 애니메이션
        'slide-in-left': {
          '0%': {
            opacity: '0',
            transform: 'translateX(-20px)'
          },
          '100%': {
            opacity: '1',
            transform: 'translateX(0)'
          },
        },
        // 오른쪽에서 들어오는 애니메이션
        'slide-in-right': {
          '0%': {
            opacity: '0',
            transform: 'translateX(20px)'
          },
          '100%': {
            opacity: '1',
            transform: 'translateX(0)'
          },
        },
      },
      animation: {
        'slide-in-left': 'slide-in-left 0.4s cubic-bezier(0.250, 0.460, 0.450, 0.940) both',
        'slide-in-right': 'slide-in-right 0.4s cubic-bezier(0.250, 0.460, 0.450, 0.940) both',
      },
    },
  },
  plugins: [],
}
