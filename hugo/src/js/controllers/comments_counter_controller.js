import $script from 'scriptjs';
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

    $script.get('https://remark42.radio-t.com/web/counter.js', () => {});
  }
}
