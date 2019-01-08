import Controller from '../base_controller';

export default class extends Controller {
  initialize() {
    super.initialize();
    window.remark_config = window.remark_config || {};
    window.remark_config.url = 'https://radio-t.com' + location.pathname;
    (function() {
      var d = document, s = d.createElement('script');
      var baseurl = 'https://remark42.radio-t.com';
      s.src = baseurl + '/web/embed.js';
      s.type = 'text/javascript';
      (d.head || d.body).appendChild(s);
    })();
  }
}
