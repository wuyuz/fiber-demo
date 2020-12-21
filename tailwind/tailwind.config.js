module.exports = {
  purge: [
    '../resource/**/*.html',
    '../resource/*.html',
    '../resource/**/*.js',
    '../resource/**/*.vue',
    '../resource/**/*.scss',
    '../resource/**/*.css',
  ],
  theme: {
    extend: {
      colors: {
        black: '#0f1c33',
      },
      margin: {
        '96': '24rem',
        '128': '32rem',
      },
    }
  },
  variants: {
    tableLayout: ['responsive', 'hover', 'focus'],
  },
  plugins: [],
}
