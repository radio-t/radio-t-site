const mix = require('laravel-mix');
const ModernizrWebpackPlugin = require('modernizr-webpack-plugin');

// mix.setPublicPath('static/build');
// mix.setResourceRoot('/build');
// Mix.manifest.name = '../../data/manifest.json';

mix.js('src/js/app.js', '.');
mix.sass('src/scss/app.scss', '.');

mix.babelConfig({
  plugins: [
    ['@babel/plugin-transform-react-jsx', {'pragma': 'h'}],
    '@babel/plugin-proposal-class-properties',
  ],
});

mix.webpackConfig({
  plugins: [
    new ModernizrWebpackPlugin(require('./.modernizr')),
  ],
});

if (!mix.inProduction()) {
  mix.setPublicPath('dev/build');
  mix.setResourceRoot('/build');

  mix.sourceMaps();
  mix.webpackConfig({devtool: 'inline-source-map'});

  mix.browserSync({
    // host: 'localhost',
    // port: 3000,
    serveStatic: ['./dev'],
    proxy: {
      target: 'localhost:1313',
      ws: true, // support websockets for hugo live-reload
    },
    files: ['dev/build/app.css', 'dev/build/app.js'],
    // watch: true,
    open: false, // don't open in browser
    ignore: ['mix-manifest.json'],
    snippetOptions: {
      rule: {
        match: /<\/head>/i,
        fn: function (snippet, match) {
          return snippet + match;
        }
      }
    },
  });
}

// mix.version();

// Full API
// mix.js(src, output);
// mix.react(src, output); <-- Identical to mix.js(), but registers React Babel compilation.
// mix.preact(src, output); <-- Identical to mix.js(), but registers Preact compilation.
// mix.coffee(src, output); <-- Identical to mix.js(), but registers CoffeeScript compilation.
// mix.ts(src, output); <-- TypeScript support. Requires tsconfig.json to exist in the same folder as webpack.mix.js
// mix.extract(vendorLibs);
// mix.sass(src, output);
// mix.less(src, output);
// mix.stylus(src, output);
// mix.postCss(src, output, [require('postcss-some-plugin')()]);
// mix.browserSync('my-site.test');
// mix.combine(files, destination);
// mix.babel(files, destination); <-- Identical to mix.combine(), but also includes Babel compilation.
// mix.copy(from, to);
// mix.copyDirectory(fromDir, toDir);
// mix.minify(file);
// mix.sourceMaps(); // Enable sourcemaps
// mix.version(); // Enable versioning.
// mix.disableNotifications();
// mix.setPublicPath('path/to/public');
// mix.setResourceRoot('prefix/for/resource/locators');
// mix.autoload({}); <-- Will be passed to Webpack's ProvidePlugin.
// mix.webpackConfig({}); <-- Override webpack.config.js, without editing the file directly.
// mix.babelConfig({}); <-- Merge extra Babel configuration (plugins, etc.) with Mix's default.
// mix.then(function () {}) <-- Will be triggered each time Webpack finishes building.
// mix.extend(name, handler) <-- Extend Mix's API with your own components.
// mix.options({
//   extractVueStyles: false, // Extract .vue component styling to file, rather than inline.
//   globalVueStyles: file, // Variables file to be imported in every component.
//   processCssUrls: true, // Process/optimize relative stylesheet url()'s. Set to false, if you don't want them touched.
//   purifyCss: false, // Remove unused CSS selectors.
//   terser: {}, // Terser-specific options. https://github.com/webpack-contrib/terser-webpack-plugin#options
//   postCss: [] // Post-CSS options: https://github.com/postcss/postcss/blob/master/docs/plugins.md
// });
