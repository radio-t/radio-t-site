if (process.env.NODE_ENV === 'production' || process.env.ENABLE_SENTRY) {
  require('./sentry');
}

require('./polyfills');
require('./stimulus');

if (process.env.NODE_ENV === 'production' || process.env.MIX_TURBO) {
  require('./turbolinks');
  require('./quicklink');
}

require('./highlight');
require('./icons');
