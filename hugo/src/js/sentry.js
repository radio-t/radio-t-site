export function initSentry() {
  const Sentry = import('@sentry/browser');
  Sentry.init({
    dsn: 'https://6571368ba3af42308da7865628a950b6@sentry.io/1467904',
  });
  return Sentry;
}