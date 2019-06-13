function getStylesheets() {
  const styles = document.querySelectorAll('link[title][rel~="stylesheet"]');
  return new Set([].slice.call(styles));
}

const styleSheets = getStylesheets();

let dark = true;

function enableStylesheet(link, enable) {
  const rel = enable ? 'stylesheet' : 'alternate stylesheet';
  link.setAttribute('rel', rel);
  link.disabled = true;
  link.disabled = !enable;
  // if (enable) {
  //   link.disabled = true;
  //   link.disabled = !enable;
  // } else {
  //   setTimeout(() => {
  //     link.disabled = !enable;
  //   }, 100);
  // }
}

function setTheme(theme) {
  console.log(styleSheets);
  getStylesheets().forEach((link) => {
    enableStylesheet(link, link.getAttribute('title') === theme);
  });
}

function getThemeName(isDark) {
  return isDark ? 'Dark' : 'Light';
}

function toggleTheme() {
  dark = !dark;
  setTheme(getThemeName(dark));
}

/**
 * Remove extra stylesheets after turbolinks render and head merge
 */
export function ensureTheme() {
  // setTheme(getThemeName(dark));
  // styleSheets.forEach((link) => {
  //   document.head.appendChild(link);
  // });
  // getStylesheets().forEach((link) => {
  //   if (!styleSheets.has(link)) {
  //     link.remove();
  //   }
  // });
}

// setTheme(getThemeName(dark));

window.setTheme = setTheme;
window.toggleTheme = toggleTheme;
