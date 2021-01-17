import Controller from '../base_controller';

const titleAttribute = 'data-theme';
export default class extends Controller {
  initialize() {
    super.initialize();

    if (!this.isSaved()) {
      window.matchMedia('(prefers-color-scheme: dark)').addListener((e) => {
        this.setTheme(e.matches ? 'dark' : 'light', false);
      });
    }
  }

  isSaved() {
    try {
      return Boolean(localStorage.getItem('theme'));
    } catch (e) {
      return false;
    }
  }

  enableStylesheet(link, enable) {
    // link.media = enable ? '' : 'none';
    if (enable) {
      link.media = '';
    } else {
      // Delay disabling to prevent FOUC
      setTimeout(() => (link.media = 'none'), 100);
    }
  }

  setTheme(theme, save = true) {
    try {
      if (save) {
        localStorage.setItem('theme', theme);
      }
    } catch (e) {
      //
    }

    const styles = [...document.querySelectorAll(`link[${titleAttribute}][rel~="stylesheet"]`)];

    styles.forEach((link) => {
      this.enableStylesheet(link, link.getAttribute(titleAttribute) === theme);
    });

    window.RADIOT_THEME = theme;

    const event = document.createEvent('Events');
    event.initEvent('theme:change', true);
    document.dispatchEvent(event);
  }

  toggle() {
    this.setTheme(window.RADIOT_THEME === 'dark' ? 'light' : 'dark');
  }
}
