module.exports = {
  content: [
    "./layouts/**/*.html",
    "./static/css/**/*.css",
    "./static/js/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        'stars-primary': '#0a192f',
        'stars-secondary': '#112240',
        'stars-accent': '#64ffda',
        'stars-gold': '#ffd700',
        'stars-text': '#ccd6f6',
        'stars-muted': '#8892b0',
      },
      keyframes: {
        'toast-in': {
          '0%': { transform: 'translateX(-50%) translateY(1rem)', opacity: 0 },
          '100%': { transform: 'translateX(-50%) translateY(0)', opacity: 1 },
        },
        'toast-out': {
          '0%': { transform: 'translateX(-50%) translateY(0)', opacity: 1 },
          '100%': { transform: 'translateX(-50%) translateY(1rem)', opacity: 0 },
        }
      },
      animation: {
        'toast-in': 'toast-in 0.3s ease-out forwards',
        'toast-out': 'toast-out 0.3s ease-in forwards'
      }
    }
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
    require('daisyui')
  ],
  daisyui: {
    themes: false,
    darkTheme: "dark",
  }
} 