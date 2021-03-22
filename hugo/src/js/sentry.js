if (process.env.NODE_ENV === 'production') {
  const captureError = async (error) => {
    console.error(error); // Log error to console for clarification
    import(
      /* webpackChunkName: "sentry" */
      /* webpackMode: "lazy" */
      '@sentry/browser'
    )
      .then((Sentry) => {
        Sentry.init({
          dsn: 'https://6571368ba3af42308da7865628a950b6@sentry.io/1467904',
        });
        Sentry.captureException(error);
      })
      .catch(() => {
        // all fails, reset window.onerror to prevent infinite loop on window.onerror
        console.error('Logging to Sentry failed', e);
        window.onerror = null;
      });
  };

  window.onerror = (message, url, line, column, error) => captureError(error);
  window.onunhandledrejection = (event) => captureError(event.reason);
}
