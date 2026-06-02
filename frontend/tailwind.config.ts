import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './composables/**/*.{js,ts}',
    './app.vue',
    './error.vue',
  ],
  theme: {
    extend: {
      colors: {
        void: 'var(--void)',
        'void-2': 'var(--void-2)',
        surface: 'var(--surface)',
        'surface-2': 'var(--surface-2)',
        'surface-solid': 'var(--surface-solid)',
        ink: 'var(--ink)',
        'ink-dim': 'var(--ink-dim)',
        'ink-faint': 'var(--ink-faint)',
        line: 'var(--line)',
        'line-strong': 'var(--line-strong)',
        primary: 'var(--primary)',
        'primary-deep': 'var(--primary-deep)',
        secondary: 'var(--secondary)',
        magenta: 'var(--magenta)',
        gold: 'var(--gold)',
        emerald: 'var(--emerald)',
        ember: 'var(--ember)',
      },
      fontFamily: {
        display: ['Cinzel', 'serif'],
        deco: ['Cinzel Decorative', 'serif'],
        body: ['Manrope', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      borderRadius: {
        sm: '8px',
        md: '14px',
        lg: '22px',
        xl: '30px',
      },
      boxShadow: {
        'glow-primary': '0 0 28px rgba(70, 232, 255, 0.4)',
        'glow-secondary': '0 0 28px rgba(183, 104, 255, 0.4)',
        panel: '0 24px 60px -28px rgba(0, 0, 0, 0.8), 0 1px 0 0 rgba(255, 255, 255, 0.04) inset',
      },
      backgroundImage: {
        arcane: 'linear-gradient(120deg, #46e8ff 0%, #b768ff 55%, #ff4fa3 100%)',
        ember: 'linear-gradient(120deg, #ffc24d, #ff6a55)',
      },
    },
  },
  plugins: [],
}
