/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./cli/client/*.{html, js}"],
  theme: {
    extend: {},
  },
  plugins: [require('@tailwindcss/typography')],
}
