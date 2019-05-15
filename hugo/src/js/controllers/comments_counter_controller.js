import $script from 'scriptjs';
import Controller from '../base_controller';
import { getUnits } from '../utils';

export default class extends Controller {
  static targets = ['number'];

  initialize() {
    super.initialize();

    const tmp = document.createElement('a');

    this.numberTargets.forEach((num) => {
      tmp.href = num.getAttribute('data-url');
      num.setAttribute('data-url', 'https://radio-t.com' + (new URL(tmp.href)).pathname);
    });

    const callback = (mutations) => {
      mutations.forEach((mutation) => {
        const num = mutation.target;
        const n = parseInt(num.innerText);
        if (isNaN(n)) return;
        num.nextElementSibling.innerHTML = getUnits(n, ['комментарий', 'комментария', 'комментариев']);
      });
    };
    const observer = new MutationObserver(callback);
    this.numberTargets.forEach((num) => {
      observer.observe(num, {
        characterData: true,
        attributes: true,
        childList: true,
        subtree: true,
      });
    });

    $script.get('https://remark42.radio-t.com/web/counter.js', () => {});
  }
}
