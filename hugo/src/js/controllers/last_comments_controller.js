import Controller from '../base_controller';

export default class extends Controller {
  initialize() {
    super.initialize();
    (function () {
      var d = document, s = d.createElement('script');
      var baseurl = 'https://remark42.radio-t.com';
      s.src = baseurl + '/web/last-comments.js';
      s.type = 'text/javascript';
      (d.head || d.body).appendChild(s);
    })();
  }
}
