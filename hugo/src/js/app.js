// TODO: fix babel-loader and use import
require('./polyfills');

let sentryInitialized = false;

async function loadSentry() {
  if (!sentryInitialized) {
    sentryInitialized = true;
    const {initSentry} = await import('./sentry');
    await initSentry();
  }
}

if (process.env.NODE_ENV !== 'production') {
  // Include here for dev, but inline for prod
  require('./inline');
}

require('./stimulus');

if (process.env.NODE_ENV === 'production' || process.env.MIX_TURBO) {
  require('./turbolinks');
  require('./quicklink');
}

// Lazy load Sentry only on error
if (process.env.NODE_ENV === 'production' || process.env.ENABLE_SENTRY) {
  window.addEventListener('error', loadSentry);
  window.addEventListener('unhandledrejection', loadSentry);
}