import Controller from '../base_controller';

const titleAttribute = 'data-theme';
let theme = window.RADIOT_THEME;

export default class extends Controller {
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

  setTheme(theme) {
    try {
      localStorage.setItem('theme', theme);
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
    theme = 'dark' === theme ? 'light' : 'dark';

    this.setTheme(theme);
  }
}
