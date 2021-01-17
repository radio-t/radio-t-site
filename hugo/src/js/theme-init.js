(function (window, themesList) {
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

  for (const t in themesList) {
    const active = theme === t;
    themesList[t].forEach((style) => {
      document.writeln(
        `<link href="${style}" rel="stylesheet" data-theme="${t}" media="${active ? '' : 'none'}">`
      );
    });
  }

  window.RADIOT_THEME = theme;
})(window, window.RADIOT_THEMES);
