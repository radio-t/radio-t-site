import { Controller } from 'stimulus';
import filter from 'lodash/filter';
import find from 'lodash/find';

export default class extends Controller {
  // static targets = ['time'];

  initialize() {
    this.fixTopics();
    this.element.classList.remove('no-js');
  }

  fixTopics() {
    for (let li of this.element.querySelectorAll('ul:first-of-type li')) {
      const timeLabel = find(li.children, (child) => child.tagName === 'EM' && child.textContent.match(/^\d+:\d+:\d+$/));
      if (!timeLabel) continue;
      // timeLabel.dataset.target = `${this.identifier}.time`;
      timeLabel.dataset.action = `click->${this.identifier}#gototime`;

      li.insertBefore(timeLabel, li.firstChild);

      while (li.childNodes.length && li.childNodes[li.childNodes.length - 1].nodeName === 'BR') {
        li.childNodes[li.childNodes.length - 1].remove();
      }
      filter(li.childNodes, (child) => {
        return child.nodeName === '#text' && child.textContent.match(/^[ \-.]+$/);
      }).forEach((node) => node.remove());
    }
  }

  gototime(e) {
    console.log(e);
  }
}
