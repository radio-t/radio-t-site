import find from 'lodash/find';

import Controller from '../base_controller';
import { composeTime, parseTime } from '../utils';

export default class extends Controller {
  connect() {
    super.connect();
    this.timeLabels();
    this.removeFirstImage();
    this.element.classList.remove('no-js');
  }

  removeFirstImage() {
    const image = this.element.querySelector('p:first-child > img:first-child');
    if (image) image.remove();
  }

  timeLabels() {
    function isEmpty(child) {
      return (
        (child.nodeName === '#text' && child.textContent.match(/^[\s\-.]+$/)) ||
        child.nodeName === 'BR'
      );
    }

    for (let li of this.element.querySelectorAll('ul:first-of-type li')) {
      let timeLabel = find(li.children, (child) => {
        return child.tagName === 'EM' && child.textContent.match(/^(\d+:)?\d+:\d+$/);
      });

      if (timeLabel) {
        timeLabel.remove();
        timeLabel.textContent = composeTime(parseTime(timeLabel.textContent));
        timeLabel.dataset.action = `click->podcast#goToTimeLabel`;
        const icon = document.createElement('i');
        icon.innerHTML = `<svg width="18" height="18" viewBox="0 0 512 512"><use xlink:href="#icon-forward-step" /></svg>`;
        timeLabel.insertBefore(icon, timeLabel.firstChild);
      } else {
        timeLabel = document.createElement('em');
      }

      // Remove empty nodes
      while (li.lastChild && isEmpty(li.lastChild)) {
        li.lastChild.remove();
      }
      if (li.childNodes && li.lastChild.nodeName === '#text') {
        li.lastChild.textContent = li.lastChild.textContent.replace(/[\s\-.]+$/, '');
      }

      // Wrap all content except time label
      const wrapper = document.createElement('div');
      wrapper.classList.add('podcast-topic-label');
      while (li.firstChild) wrapper.appendChild(li.firstChild);

      // Put into dom
      li.appendChild(timeLabel);
      li.appendChild(wrapper);
    }
  }
}
