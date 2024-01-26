/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["web/view/**/*.templ"],
  theme: {
    extend: {
      fontFamily: {
        'title': ['TitleRegular', 'sans-serif'],
        'text': ['TextRegular', 'sans-serif'],
        'text-bold': ['TextBold', 'sans-serif'],
        'italic': ['TextItalic', 'sans-serif'],
      },
      borderColor: {
        'accent': 'var(--accent-color)',
      },
      colors: {
        'accent': {
          DEFAULT: 'var(--accent-color)',
          '5': 'rgba(8, 76, 223, 0.05)',
          '10': 'rgba(8, 76, 223, 0.10)',
          '15': 'rgba(8, 76, 223, 0.15)',
          '30': 'rgba(8, 76, 223, 0.3)',
          '80': 'rgba(8, 76, 223, 0.8)',
        },
      }
    },
  },
  plugins: [],
}

