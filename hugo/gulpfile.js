const gulp = require('gulp');
const favicons = require('gulp-favicons');

gulp.task('favicons', function () {
  return gulp.src('./favicon.png')
    .pipe(favicons({
      appName: 'Радио-Т Подкаст',
      appDescription: 'Подкаст выходного дня - импровизации на темы высоких технологий',
      // developerName: '',
      // developerURL: '',
      background: 'transparent',
      path: '/',
      // url: '',
      // display: 'standalone',
      // orientation: 'portrait',
      // start_url: '/?homescreen=1',
      // version: 1.0,
      logging: true,
      online: false,
      html: '../layouts/partials/favicons.html',
      pipeHTML: true,
      replace: true,
      icons: {
        // Platform Options:
        // - offset - offset in percentage
        // - background:
        //   * false - use default
        //   * true - force use default, e.g. set background for Android icons
        //   * color - set background for the specified icons
        //   * mask - apply mask in order to create circle icon (applied by default for firefox). `boolean`
        //   * overlayGlow - apply glow effect after mask has been applied (applied by default for firefox). `boolean`
        //   * overlayShadow - apply drop shadow after mask has been applied .`boolean`
        //
        favicons: true,              // Create regular favicons. `boolean` or `{ offset, background, mask, overlayGlow, overlayShadow }`
        android: false,              // Create Android homescreen icon. `boolean` or `{ offset, background, mask, overlayGlow, overlayShadow }`
        appleIcon: false,            // Create Apple touch icons. `boolean` or `{ offset, background, mask, overlayGlow, overlayShadow }`
        appleStartup: false,         // Create Apple startup images. `boolean` or `{ offset, background, mask, overlayGlow, overlayShadow }`
        coast: false,                // Create Opera Coast icon. `boolean` or `{ offset, background, mask, overlayGlow, overlayShadow }`
        firefox: false,              // Create Firefox OS icons. `boolean` or `{ offset, background, mask, overlayGlow, overlayShadow }`
        windows: false,              // Create Windows 8 tile icons. `boolean` or `{ offset, background, mask, overlayGlow, overlayShadow }`
        yandex: false,               // Create Yandex browser icon. `boolean` or `{ offset, background, mask, overlayGlow, overlayShadow }`
      },
    }))
    .on('error', require('fancy-log'))
    .pipe(gulp.dest('./static/'));
});
