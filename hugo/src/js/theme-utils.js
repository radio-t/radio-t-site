const titleAttribute = 'data-theme';
let theme = window.RADIOT_THEME;

const isDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches
const isLightMode = window.matchMedia("(prefers-color-scheme: light)").matches
const isNotSpecified = window.matchMedia("(prefers-color-scheme: no-preference)").matches
const hasNoSupport = !isDarkMode && !isLightMode && !isNotSpecified;

window.matchMedia("(prefers-color-scheme: dark)").addListener(e => e.matches && setTheme("dark"));
window.matchMedia("(prefers-color-scheme: light)").addListener(e => e.matches && setTheme("light"));

export function getStylesheets() {
    const styles = document.querySelectorAll(`link[${titleAttribute}][rel~="stylesheet"]`);
    return [].slice.call(styles);
}

export function enableStylesheet(link, enable) {
    // link.media = enable ? '' : 'none';
    if (enable) {
        link.media = '';
    } else {
        // Delay disabling to prevent FOUC
        setTimeout(() => link.media = 'none', 100);
    }
}

export function setTheme(theme) {
    // theme = 'dark' === theme ? 'light' : 'dark';
    try {
        localStorage.setItem('theme', theme);
    } catch (e) {
        //
    }

    getStylesheets().forEach((link) => {
        enableStylesheet(link, link.getAttribute(titleAttribute) === theme);
    });

    window.RADIOT_THEME = theme;

    const event = document.createEvent('Events');
    event.initEvent('theme:change', true);
    document.dispatchEvent(event);
}