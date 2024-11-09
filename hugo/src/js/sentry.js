let Sentry;

/** Lazy load Sentry on error */
async function handleException(errorEvent) {
  try {
    if (Sentry === undefined) {
      Sentry = await import('@sentry/browser');
      Sentry.init({
        dsn: 'https://86c7b8de1ad3cf69978fdf409a776f28@o510231.ingest.us.sentry.io/4508265848897536',
        enabled: process.env.NODE_ENV === 'production' || process.env.ENABLE_SENTRY,
      });
    }
    Sentry.captureException(errorEvent);
  } catch (e) {
    console.error('Logging to Sentry failed', e);
  }
}

window.addEventListener('error', handleException);
window.addEventListener('unhandledrejection', handleException);
