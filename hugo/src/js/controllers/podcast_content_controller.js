import throttle from 'lodash/throttle';
import find from 'lodash/find';

import Controller from '../base_controller';
import { composeTime, parseTime, getLocalStorage } from '../utils';

export default class extends Controller {
  topics = [];
  activeTopic = null;
  podcastNumber = null;

  initialize() {
    super.initialize();

    const postPodcast = this.element.closest('.post-podcast');
    const podcastTitleNumber = postPodcast?.querySelector('.podcast-title-number');
    this.podcastNumber = podcastTitleNumber?.innerText?.trim() || null;

    if (!this.podcastNumber) {
      return;
    }

    this.subscribe(
      `playing-progress-${this.podcastNumber}`,
      throttle(this.updateActiveTopic.bind(this), 1000)
    );
  }

  connect() {
    super.connect();
    this.timeLabels();
    this.removeFirstImage();
    this.setInitialActiveTopic();
    this.element.classList.remove('no-js');
  }

  setInitialActiveTopic() {
    if (!this.podcastNumber) {
      return;
    }

    const podcasts = getLocalStorage(`podcasts`) || {};
    if (podcasts[this.podcastNumber]) {
      const { currentTime } = podcasts[this.podcastNumber];
      this.updateActiveTopic({ currentTime });
    }
  }

  findTopicForTime(time) {
    let result = null;

    for (const [topicTime, el] of this.topics) {
      if (topicTime <= time) {
        result = el;
      } else {
        break;
      }
    }

    return result;
  }

  updateActiveTopic({ currentTime = 0 }) {
    const currentTopic = this.findTopicForTime(currentTime);

    if (!currentTopic || currentTopic === this.activeTopic) {
      return;
    }

    if (this.activeTopic) {
      this.activeTopic.classList.remove('active');
    }

    this.activeTopic = currentTopic;
    currentTopic.classList.add('active');
  }

  removeFirstImage() {
    const image = this.element.querySelector('p:first-child > img:first-child');
    if (image) image.remove();
  }

  timeLabels() {
    this.topics = [];
    this.activeTopic = null;

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
        const time = parseTime(timeLabel.textContent);
        this.topics.push([time, li]);
        timeLabel.textContent = composeTime(time);
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
      if (li.lastChild && li.lastChild.nodeName === '#text') {
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
