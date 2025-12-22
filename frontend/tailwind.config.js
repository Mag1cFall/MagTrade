/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // 极简黑白灰基调
        background: '#050505', // 深邃黑背景
        surface: '#121212',    // 卡片/表面颜色
        'surface-light': '#1E1E1E', // 略亮的表面
        border: '#2A2A2A',     // 极细边框
        
        // 文字颜色
        primary: '#FFFFFF',    // 主标题
        secondary: '#A1A1AA',  // 次要文字/正文
        tertiary: '#52525B',   // 辅助文字/占位符

        // 强调色 - 赤焰红 (仅用于行动点和紧迫感)
        accent: '#FF3B30',
        'accent-hover': '#FF2D20',
        'accent-dim': 'rgba(255, 59, 48, 0.1)', // 红色微光背景
      },
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'ui-monospace', 'monospace'], // 用于数字和代码
      },
      animation: {
        'pulse-fast': 'pulse 1.5s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'fade-in': 'fadeIn 0.5s ease-out',
        'slide-up': 'slideUp 0.5s ease-out',
        'pulse-red': 'pulse-red 2s infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        'pulse-red': {
          '0%, 100%': { boxShadow: '0 0 0 0 rgba(255, 59, 48, 0)' },
          '50%': { boxShadow: '0 0 15px 0 rgba(255, 59, 48, 0.3)' },
        },
      },
    },
  },
  plugins: [],
}