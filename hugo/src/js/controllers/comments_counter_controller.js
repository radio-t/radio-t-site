import Controller from '../base_controller';

export default class extends Controller {
  static targets = ['number'];

  initialize() {
    super.initialize();

    const tmp = document.createElement('a');

    this.numberTargets.forEach((s) => {
      tmp.href = s.getAttribute('data-url');
      s.setAttribute('data-url', 'https://radio-t.com' + (new URL(tmp.href)).pathname);
    });

    this.runScript();
  }

  runScript() {
    var d = document, s = d.createElement('script');
    var baseurl = 'https://remark42.radio-t.com';
    s.src = baseurl + '/web/counter.js';
    s.type = 'text/javascript';
    (d.head || d.body).appendChild(s);
  }
}
