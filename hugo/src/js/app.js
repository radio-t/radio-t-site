// TODO: fix babel-loader and use import
require('./polyfills');

if (process.env.NODE_ENV !== 'production') {
  // Include here for dev, but inline for prod
  require('./inline');
}

require('./stimulus');

if (process.env.NODE_ENV === 'production' || process.env.MIX_TURBO) {
  require('./turbolinks');
  require('./quicklink');
}

if (process.env.NODE_ENV === 'production' || process.env.ENABLE_SENTRY) {
  require('./sentry');
}
