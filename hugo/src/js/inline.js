((window, themes) => {
  // Insert CSS
  const theme = (() => {
    try {
      const themeFromStore = localStorage.getItem('theme');

      if (['dark', 'light'].includes(themeFromStore)) {
        return themeFromStore;
      }

      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    } catch (e) {
      return 'light';
    }
  })();

  for (const t in themes) {
    const active = theme === t;
    themes[t].forEach((style) => {
      document.writeln(
        `<link href="${style}" rel="stylesheet" data-theme="${t}" media="${active ? '' : 'none'}">`
      );
    });
  }

  window.RADIOT_THEME = theme;

  // SVG SPRITE
  const container = document.createElement('div');

  container.className = 'd-none';

  function insert() {
    document.body.prepend(container);
  }

  fetch('/build/images/icons-sprite.svg')
    .then((r) => r.text())
    .then((data) => {
      container.innerHTML = data;

      if (document.body) {
        insert();
        return;
      }

      document.onreadystatechange = function () {
        if (document.readyState === 'interactive') {
          insert();
        }
      };
    });
})(window, window.RADIOT_THEMES);
