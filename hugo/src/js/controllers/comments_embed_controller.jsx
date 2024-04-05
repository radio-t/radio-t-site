import Controller from '../base_controller';

export default class extends Controller {
  get theme() {
    return window.RADIOT_THEME === 'dark' ? 'dark' : 'light';
  }
  initialize() {
    super.initialize();
    window.remark_config.url = `https://radio-t.com${location.pathname}`;
    window.remark_config.theme = this.theme;
    if (window.REMARK42) {
      window.REMARK42.destroy();
      window.REMARK42.createInstance(window.remark_config);
    }
    this.changeTheme = this.changeTheme.bind(this);
  }

  connect() {
    super.connect();
    document.addEventListener('theme:change', this.changeTheme);
  }

  disconnect() {
    super.disconnect();
    document.removeEventListener('theme:change', this.changeTheme);
  }

  changeTheme() {
    window.remark_config.theme = this.theme;
    window.REMARK42.changeTheme(this.theme);
  }
}
