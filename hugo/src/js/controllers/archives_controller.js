import { Controller } from 'stimulus';
import uniq from 'lodash/uniq';
import find from 'lodash/find';
import findLast from 'lodash/findLast';

export default class extends Controller {
  static targets = ['range', 'list'];

  ranges;

  initialize() {
    this.ranges = this.listTargets.map((list) => {
      const latest = find(list.children, (post) => post.querySelector('.podcast-title-number'));
      const earliest = findLast(list.children, (post) => post.querySelector('.podcast-title-number'));
      return uniq([earliest, latest]).map((post) => {
        if (post) return post.querySelector('.podcast-title-number').textContent;
      }).filter(s => s).join(' – ');
    });

    this.rangeTargets.forEach((target, index) => target.textContent = this.ranges[index]);
  }

  connect() {
    this.rangeTargets.forEach((target, index) => {
      target.style.transitionDelay = `${350 + index * 50}ms`;
      target.style.transitionDuration = `200ms`;
      target.style.transitionProperty = 'opacity';
      target.style.opacity = 1;
    });
  }

  disconnect() {
    this.rangeTargets.forEach((target) => {
      target.style.transitionDelay = `0ms`;
      target.style.opacity = 0;
    });
  }
}
