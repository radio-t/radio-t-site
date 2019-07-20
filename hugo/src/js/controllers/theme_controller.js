import Controller from '../base_controller';

const titleAttribute = 'data-theme';

export default class extends Controller {
  initialize() {
    super.initialize();

    window.matchMedia("(prefers-color-scheme: dark)").addListener(e => e.matches && this.onMatchMedia("dark"));
    window.matchMedia("(prefers-color-scheme: light)").addListener(e => e.matches && this.onMatchMedia("light"));
  }

  toggle() {
    this.toggleTheme();
  }

  getStylesheets() {
    const styles = document.querySelectorAll(`link[${titleAttribute}][rel~="stylesheet"]`);
    return [].slice.call(styles);
  }

  enableStylesheet(link, enable) {
    // link.media = enable ? '' : 'none';
    if (enable) {
      link.media = '';
    } else {
      // Delay disabling to prevent FOUC
      setTimeout(() => link.media = 'none', 100);
    }
  }

  onMatchMedia(theme) {
    if (localStorage.getItem('theme')) return;
    this.setTheme(theme, false);
  }

  setTheme(theme, save = true) {
    try {
      if (save) {
        localStorage.setItem('theme', theme);
      }
    } catch (e) {
      //
    }

    this.getStylesheets().forEach((link) => {
      this.enableStylesheet(link, link.getAttribute(titleAttribute) === theme);
    });

    window.RADIOT_THEME = theme;

    const event = document.createEvent('Events');
    event.initEvent('theme:change', true);
    document.dispatchEvent(event);
  }

  toggleTheme() {
    let theme = window.RADIOT_THEME;
    theme = 'dark' === theme ? 'light' : 'dark';

    this.setTheme(theme)
  }
}
