import domready from 'domready';
import quicklink from 'quicklink/dist/quicklink.mjs';

function start(options = {}) {
  options = {
    ignores: [
      (uri, elem) => String(elem.getAttribute('href'))[0] === '#',
      (uri, elem) => elem.matches('[noprefetch]') || elem.closest('[noprefetch]'),
    ],
    ...options,
  };
  quicklink(options);
}

domready(start);

document.addEventListener('turbolinks:load', start);

document.addEventListener('quicklink', (e) => start(e.detail || {}));
