(function () {
  function getSystemTheme() {
    const isDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches
    const isLightMode = window.matchMedia("(prefers-color-scheme: light)").matches
    const isNotSpecified = window.matchMedia("(prefers-color-scheme: no-preference)").matches
    const hasNoSupport = !isDarkMode && !isLightMode && !isNotSpecified;

    return hasNoSupport ? undefined
      : isDarkMode ? 'dark' : 'light'
  }

  function getPreferredTheme() {
    try {
      let theme = localStorage.getItem('theme');
      return theme
        ? 'dark' === theme ? 'dark' : 'light'
        : undefined;
    } catch (e) {
      //
    }

    return undefined;
  }

  window.RADIOT_THEME = getTheme();

  for (const t in window.RADIOT_THEMES) {
    const active = window.RADIOT_THEME === t;
    window.RADIOT_THEMES[t].forEach(function (style) {
      const tag = `<link href="${style}" rel="stylesheet" data-theme="${t}" media="${active ? '' : 'none'}">`;
      document.writeln(tag);
    });
  }
}());
