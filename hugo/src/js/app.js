import './polyfills';

if (process.env.NODE_ENV !== 'production') {
  // Include here for dev, but inline for prod
  import './theme-init';
}

import './stimulus';

if (process.env.NODE_ENV === 'production' || process.env.MIX_TURBO) {
  import './turbolinks';
  import './quicklink';
}

if (process.env.NODE_ENV === 'production' || process.env.ENABLE_SENTRY) {
  import './sentry';
}
