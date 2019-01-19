const gulp = require('gulp')
const favicons = require('gulp-favicons')

gulp.task('favicons', function () {
  return gulp.src('./favicon.png')
    .pipe(favicons({
      appName: '',
      appDescription: '',
      // developerName: '',
      // developerURL: '',
      background: '#fff',
      path: '/favicons/',
      // url: '',
      // display: 'standalone',
      // orientation: 'portrait',
      // start_url: '/?homescreen=1',
      // version: 1.0,
      logging: true,
      online: false,
      html: '../favicons.html',
      pipeHTML: true,
      replace: true,
    }))
    .on('error', require('fancy-log'))
    .pipe(gulp.dest('./favicons/'))
})
