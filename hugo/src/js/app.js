require('./polyfills');

if (process.env.NODE_ENV === 'production' || process.env.ENABLE_SENTRY) {
  require('./sentry');
}

if (process.env.NODE_ENV !== 'production') {
  // Include here for dev, but inline for prod
  require('./theme-init');
}

require('./stimulus');

if (process.env.NODE_ENV === 'production' || process.env.MIX_TURBO) {
  require('./turbolinks');
  require('./quicklink');
}

require('./highlight');
require('./icons');
