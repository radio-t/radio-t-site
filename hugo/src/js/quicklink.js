import domready from 'domready';
import quicklink from 'quicklink/dist/quicklink.mjs';

domready(() => quicklink({
  ignores: [
    (uri, elem) => String(elem.getAttribute('href'))[0] === '#',
    (uri, elem) => elem.matches('[noprefetch]') || elem.closest('[noprefetch]'),
  ],
}));
