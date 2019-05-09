import Controller from '../base_controller';
import uniq from 'lodash/uniq';
import find from 'lodash/find';
import findLast from 'lodash/findLast';
import SmoothScroll from 'smooth-scroll/dist/smooth-scroll';

export default class extends Controller {
  static targets = ['range', 'list'];

  ranges;
  scroll;

  initialize() {
    super.initialize();

    if (this.data.get('initialized')) return;
    this.data.set('initialized', '1');

    this.ranges = this.listTargets.map((list) => {
      const latest = find(list.children, (post) => post.querySelector('.podcast-title-number'));
      const earliest = findLast(list.children, (post) => post.querySelector('.podcast-title-number'));
      return uniq([latest, earliest]).map((post) => {
        if (post) return post.querySelector('.podcast-title-number').textContent;
      }).filter(s => s).join(' – ');
    });

    this.rangeTargets.forEach((target, index) => target.textContent = this.ranges[index]);

    this.scroll = new SmoothScroll('a[href*="#"]', {
      speed: 100,
      speedAsDuration: false,
      durationMax: 400,
      easing: 'easeOutCubic',
      updateURL: false,
      popstate: false,
    });
  }
}
