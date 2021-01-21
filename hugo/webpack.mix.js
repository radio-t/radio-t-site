const fs = require('fs');
const mix = require('laravel-mix');
const babel = require('@babel/core');
const ModernizrWebpackPlugin = require('modernizr-webpack-plugin');
const nodeSass = require('node-sass');

mix.disableNotifications();

mix
  .webpackConfig({
    resolve: {
      alias: {
        'react': 'preact/compat',
        'react-dom': 'preact/compat',
      },
    },
  })
  .ts('src/js/app.js', '.')
  .version();

['app', 'vendor'].forEach((style) => {
  mix.sass(`src/scss/${style}.scss`, '.', { implementation: nodeSass });
  mix.sass(`src/scss/${style}-dark.scss`, '.', { implementation: nodeSass });
});

mix.webpackConfig({ plugins: [new ModernizrWebpackPlugin(require('./.modernizr'))] });

if (process.env.ANALYZE) {
  const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');

  mix.webpackConfig({ plugins: [new BundleAnalyzerPlugin()] });
}

if (mix.inProduction()) {
  Mix.manifest.name = '../../data/manifest.json'; // eslint-disable-line no-undef
  mix.setPublicPath('static/build');
  mix.setResourceRoot('/build');
  mix.extract();
  mix.version(['static/build/modernizr-bundle.js']);
  mix.then(() => {
    const { code } = babel.transformFileSync('src/js/inline.js', {
      minified: true,
      presets: [
        [
          '@babel/preset-env',
          {
            modules: false,
            forceAllTransforms: true,
            useBuiltIns: false,
          },
        ],
      ],
    });
    fs.writeFileSync('static/build/inline.js', code);
  });
} else {
  mix.setPublicPath('dev');

  mix.sourceMaps();
  mix.webpackConfig({ devtool: 'inline-source-map' });

  mix.browserSync({
    host: process.env.DEV_HOST || 'localhost',
    port: process.env.DEV_PORT || 3000,
    serveStatic: ['./dev'],
    proxy: {
      target: `localhost:${process.env.HUGO_PORT || 1313}`,
      ws: true, // support websockets for hugo live-reload
    },
    // watch: true,
    // watch specific files
    files: ['dev/*.css', 'dev/app.js'],
    ghostMode: false, // disable Clicks, Scrolls & Form inputs on any device will be mirrored to all others
    open: false, // don't open in browser
    ignore: ['mix-manifest.json'],
    // to work with turbolinks
    snippetOptions: {
      rule: {
        match: /<\/head>/i,
        fn: function (snippet, match) {
          return snippet + match;
        },
      },
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
