// TODO: fix babel-loader and use import
import './polyfills';
import './stimulus';
// import './sentry';

if (process.env.NODE_ENV !== 'production') {
  // Include here for dev, but inline for prod
  require('./theme-init');
}

if (process.env.NODE_ENV === 'production' || process.env.MIX_TURBO) {
  require('./turbolinks');
  require('./quicklink');
}
