/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["web/view/**/*.templ"],
  theme: {
    extend: {
      fontFamily: {
        'title': ['TitleRegular', 'sans-serif'],
        'text': ['TextRegular', 'sans-serif'],
      }
    },
  },
  plugins: [],
}

