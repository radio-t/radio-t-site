import Controller from '../base_controller';
import filter from 'lodash/filter';
import find from 'lodash/find';
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
    for (let li of this.element.querySelectorAll('ul:first-of-type li')) {
      let timeLabel = find(li.children, child => {
        return child.tagName === 'EM' && child.textContent.match(/^(\d+:)?\d+:\d+$/);
      });

      if (timeLabel) {
        timeLabel.remove();
        timeLabel.textContent = composeTime(parseTime(timeLabel.textContent));
        timeLabel.dataset.action = `click->podcast#goToTimeLabel`;
        const icon = document.createElement('i');
        icon.className = 'fas fa-step-forward fa-fw';
        timeLabel.prepend(icon);
      } else {
        timeLabel = document.createElement('EM');
      }

      // Remove empty nodes
      filter(li.childNodes, (child) => {
        return child.nodeName === '#text' && child.textContent.match(/^[\s\-\.]+$/);
      }).forEach((node) => node.remove());
      console.log(li.childNodes);
      while (li.childNodes.length && li.childNodes[li.childNodes.length - 1].nodeName === 'BR') {
        li.childNodes[li.childNodes.length - 1].remove();
      }
      if (li.childNodes.length && li.childNodes[li.childNodes.length - 1].nodeName === '#text') {
        li.childNodes[li.childNodes.length - 1].textContent = li.childNodes[li.childNodes.length - 1].textContent.replace(/[ \-.]+$/, '');
      }

      // Wrap all content except time label
      const wrapper = document.createElement('div');
      wrapper.classList.add('podcast-topic-label');
      while (li.firstChild) wrapper.append(li.firstChild);

      // Put into dom
      li.append(timeLabel);
      li.append(wrapper);
    }
  }
}
