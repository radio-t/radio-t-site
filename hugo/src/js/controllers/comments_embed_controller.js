import $script from 'scriptjs';
import Controller from '../base_controller';

export default class extends Controller {
  initialize() {
    super.initialize();
    window.remark_config = window.remark_config || {};
    window.remark_config.url = 'https://radio-t.com' + location.pathname;
    $script.get('https://remark42.radio-t.com/web/embed.js', () => {});
  }
}
