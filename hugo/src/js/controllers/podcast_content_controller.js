import { Controller } from 'stimulus';
import filter from 'lodash/filter';
import find from 'lodash/find';

export default class extends Controller {
  initialize() {
    this.timeLabels();
    this.element.classList.remove('no-js');
  }

  timeLabels() {
    for (let li of this.element.querySelectorAll('ul:first-of-type li')) {
      let timeLabel = find(li.children, (child) => child.tagName === 'EM' && child.textContent.match(/^\d+:\d+:\d+$/));

      if (timeLabel) {
        timeLabel.remove();
        timeLabel.dataset.action = `click->${this.identifier}#jumpTime`;
        // timeLabel.dataset.target = `${this.identifier}.time`;
        const icon = document.createElement('i');
        icon.className = 'fas fa-step-forward fa-fw';
        timeLabel.prepend(icon);
      } else {
        timeLabel = document.createElement('EM');
      }

      // Remove empty nodes
      while (li.childNodes.length && li.childNodes[li.childNodes.length - 1].nodeName === 'BR') {
        li.childNodes[li.childNodes.length - 1].remove();
      }
      filter(li.childNodes, (child) => {
        return child.nodeName === '#text' && child.textContent.match(/^[ \-.]+$/);
      }).forEach((node) => node.remove());

      const wrapper = document.createElement('div');
      wrapper.classList.add('podcast-topic-label');
      while (li.firstChild) wrapper.append(li.firstChild);

      li.append(timeLabel);
      li.append(wrapper);
    }
  }

  jumpTime(e) {
    console.log(`Jump to time ${e.target.textContent}`);
  }
}
