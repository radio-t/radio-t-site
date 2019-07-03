(function () {
  function getPreferredTheme() {
    try {
      return 'dark' === localStorage.getItem('theme') ? 'dark' : 'light';
    } catch (e) {
      //
    }

    return 'light';
  }

  window.RADIOT_THEME = getPreferredTheme();

  for (const t in window.RADIOT_THEMES) {
    const active = window.RADIOT_THEME === t;
    window.RADIOT_THEMES[t].forEach(function (style) {
      const tag = `<link href="${style}" rel="stylesheet" data-theme="${t}" media="${active ? '' : 'none'}">`;
      document.writeln(tag);
    });
  }
}());
