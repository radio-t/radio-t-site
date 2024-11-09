const fs = require('fs');
const glob = require('glob');
const mix = require('laravel-mix');
const babel = require('@babel/core');
const PurgeCSSPlugin = require('purgecss-webpack-plugin');
const purgecssHtml = require('purgecss-from-html');
const process = require('process');

const shouldMinify = process.env.DO_NOT_MINIFY !== 'true';

mix.disableNotifications();

mix
  .webpackConfig({
    resolve: {
      alias: {
        'react': 'preact/compat',
        'react-dom': 'preact/compat',
      },
      fallback: {
        "buffer": require.resolve("buffer/")
      }
    },
  })
  .ts('src/js/app.js', '.')
  .extract(['@sentry/browser'], 'vendor~sentry.js')
  .version();


// Process CSS with conditional minification
['app', 'vendor'].forEach((style) => {
  mix
    .sass(`src/scss/${style}.scss`, '.')
    .options({ postCss: shouldMinify ? [require('cssnano')] : [] }); // Conditionally include cssnano
  mix
    .sass(`src/scss/${style}-dark.scss`, '.')
    .options({ postCss: shouldMinify ? [require('cssnano')] : [] }); // Conditionally include cssnano
});

mix.webpackConfig({
  plugins: [
    new PurgeCSSPlugin({
      paths: [
        ...glob.sync('layouts/**/*.html', { nodir: true }),
        ...glob.sync('src/**/*.{js,ts,jsx,tsx}', { nodir: true }),
      ],
      safelist: () => ({
        deep: [
          /is-online/,
          /has-audio/,
          /post-podcast-content/,
          /fa-step-forward/,
          /sidebar-open/,
          /comments-counter-avatars-item/,
          /loaded/,
          /highlight/,
          /language-/,
          /code/,
          /pre/,
        ],
      }),
      extractors: [
        {
          extensions: ['html'],
          extractor: purgecssHtml,
        },
        {
          extensions: ['js'],
          extractor(content) {
            const regexStr = "classList.\\w+.\\('(.*)'";
            const globalRegex = new RegExp(regexStr, 'gm');
            const localRegex = new RegExp(regexStr);
            const match = content.match(globalRegex);
            if (match === null) {
              return [];
            }
            return { classes: match.map((s) => s.match(localRegex)[1]) };
          },
        },
      ],
    }),
  ]
});

if (process.env.ANALYZE) {
  const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');

  mix.webpackConfig({ plugins: [new BundleAnalyzerPlugin()] });
}

// Production and Development Paths
if (mix.inProduction()) {
  Mix.manifest.name = '../../data/manifest.json'; // eslint-disable-line no-undef
  mix.setPublicPath('static/build');
  mix.setResourceRoot('/build');
  mix.extract();
  mix.then(() => {
    const { code } = babel.transformFileSync('src/js/inline.js', {
      minified: shouldMinify, // Conditionally set minification for JS
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
  mix.copy('src/images/icons-sprite.svg', 'static/build/images/icons-sprite.svg');
} else {
  mix.setPublicPath('dev');
  mix.copy('src/images/icons-sprite.svg', 'dev/build/images/icons-sprite.svg');

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
mix.webpackConfig({
  output: {
    publicPath: '/build',
  },
  optimization: {
    splitChunks: 'all',
  },
})