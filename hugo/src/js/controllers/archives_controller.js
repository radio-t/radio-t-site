import { Controller } from 'stimulus';
import uniq from 'lodash/uniq';
import find from 'lodash/find';
import findLast from 'lodash/findLast';

export default class extends Controller {
  static targets = ['range', 'list'];

  initialize() {
    const ranges = this.listTargets.map((list) => {
      const latest = find(list.children, (post) => post.querySelector('.podcast-title-number'));
      const earliest = findLast(list.children, (post) => post.querySelector('.podcast-title-number'));
      return uniq([earliest, latest]).map((post) => {
        if (post) return post.querySelector('.podcast-title-number').textContent;
      }).filter(s => s).join(' – ');
    });

    ranges.forEach((range, index) => {
      this.rangeTargets[index].textContent = range;
      this.rangeTargets[index].style.transitionDelay = `${index * 50}ms`;
      this.rangeTargets[index].classList.add('in');
    });
  }
}
