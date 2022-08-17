import Controller from '../base_controller';

export default class extends Controller {
  get theme() {
    return window.RADIOT_THEME === 'dark' ? 'dark' : 'light';
  }
  initialize() {
    super.initialize();
  }

  connect() {
    super.connect();
    document.addEventListener('theme:change', this.changeTheme);
  }

  render = () => {
    window.remark_config.url = `https://radio-t.com${location.pathname}`;
    window.remark_config.page_title = document.title;
    window.remark_config.theme = this.theme;
  };

  disconnect() {
    super.disconnect();
    document.removeEventListener('theme:change', this.changeTheme);
  }

  changeTheme() {
    window.remark_config.theme = this.theme;
    window.REMARK42.changeTheme(this.theme);
  }
}
