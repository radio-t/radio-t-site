import quicklink from 'quicklink/dist/quicklink.mjs';

function start(options = {}) {
  quicklink({
    ignores: [
      (_, elem) => String(elem.getAttribute('href'))[0] === '#',
      (_, elem) => elem.matches('[noprefetch]') || elem.closest('[noprefetch]'),
    ],
    ...options,
  });
}

document.addEventListener('turbolinks:load', start);
document.addEventListener('quicklink', (e) => start(e.detail || {}));

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', start);
} else {
  start();
}