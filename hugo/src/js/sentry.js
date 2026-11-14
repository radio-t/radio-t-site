let sentryLoading;

/** Unwrap the actual error from an ErrorEvent or PromiseRejectionEvent */
function getError(event) {
  if ('reason' in event) return event.reason;
  return event.error || event.message;
}

/**
 * Lazy load Sentry on the first error. After init, Sentry's own global handlers
 * take over, so only errors raised before that are reported from here.
 */
async function handleException(event) {
  try {
    if (sentryLoading === undefined) {
      sentryLoading = import('@sentry/browser').then((Sentry) => {
        Sentry.init({
          dsn: 'https://86c7b8de1ad3cf69978fdf409a776f28@o510231.ingest.us.sentry.io/4508265848897536',
          enabled: process.env.NODE_ENV === 'production' || process.env.ENABLE_SENTRY,
          ignoreErrors: [
            'Network Error',
            // Injected by a browser extension, not our code
            'MtPopUpList is not defined',
            // play() interrupted by pause() or a new load; harmless
            /The play\(\) request was interrupted/,
          ],
        });
        window.removeEventListener('error', handleException);
        window.removeEventListener('unhandledrejection', handleException);
        return Sentry;
      });
    }
    const Sentry = await sentryLoading;
    Sentry.captureException(getError(event));
  } catch (e) {
    console.error('Logging to Sentry failed', e);
  }
}

window.addEventListener('error', handleException);
window.addEventListener('unhandledrejection', handleException);
