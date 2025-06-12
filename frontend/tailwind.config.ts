import type { Config } from 'tailwindcss'

export default {
  theme: {
    fontSize: {
      'sm': ['0.625rem', { lineHeight: '0.875rem' }],
      'base': ['0.75rem', { lineHeight: '1rem' }],
      'lg': ['0.875rem', { lineHeight: '1.25rem' }],
      'xl': ['1.375rem', { lineHeight: '1.75rem' }],
      '2xl': ['2rem', { lineHeight: '2.5rem' }],
    },
    colors: {
      "dark-blue": "#323C4C",
      "dark-blue-80": "#5B6370",
      "dark-blue-55": "#8E949D",
      "dark-blue-35": "#B7BBC0",
      "dark-blue-12": "#E0E2E4",
      "dark-blue-hover": "#F5F6F8",
      "white": "#FFFFFF",
      "light-blue": "#B1EAFC",
      "light-blue-60": "#D0F2FD",
      "light-blue-40": "#E0F7FE",
      "light-blue-20": "#EFFBFE",
      "bg-gray": "#F7F8FA",
      "error": "#F87979",
      "warning": "#FFD260",
      "text": "#5B6370"
    },
    extend: {
      fontFamily: {
        notoSans: ['var(--font-noto-sans)', 'sans-serif'],
      },
      fontWeight: {
        light: 300,
        normal: 400,
        medium: 500,
        bold: 700,
      },
    },
  },
  plugins: [],
} satisfies Config;
