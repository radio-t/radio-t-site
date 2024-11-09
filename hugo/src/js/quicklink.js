import { listen } from 'quicklink';

window.addEventListener('load', () => {
  listen({
    ignores: [
      (_, elem) => String(elem.getAttribute('href'))[0] === '#',
      (_, elem) => elem.matches('[noprefetch]') || elem.closest('[noprefetch]'),
    ],
  });
});
