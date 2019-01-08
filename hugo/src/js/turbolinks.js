import Turbolinks from 'turbolinks';

Turbolinks.start();

// open external links in new tab
document.addEventListener('turbolinks:load', () => {
  [].forEach.call(document.links, link => {
    if (link.hostname !== window.location.hostname) {
      link.target = '_blank';
    }
  });
});
