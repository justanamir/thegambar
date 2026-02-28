module.exports = {
  content: ["./web/templates/**/*.html"],
  theme: {
    extend: {
      colors: {
        amber:   '#FFBF00',
        gold:    '#F2CF7E',
        yellow:  '#FFE642',
        orange:  '#FF7900',
        ink:     '#1C1400',
        cream:   '#FFFDF5',
        warm:    '#FFF3D6',
      },
      fontFamily: {
        display: ['Georgia', 'Cambria', 'serif'],
        body:    ['system-ui', 'sans-serif'],
      }
    }
  },
  plugins: [],
}