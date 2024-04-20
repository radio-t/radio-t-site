import { throttle } from 'lodash';
import Controller from '../base_controller';
import { composeTime, getLocalStorage } from '../utils';

export default class extends Controller {
  static targets = ['bar', 'number', 'duration', 'progress'];

  initialize() {
    super.initialize();
    this.subscribe(
      `playing-progress-${this.numberTarget.innerText}`,
      throttle(this.renderProgress.bind(this), 1000)
    );
  }

  connect() {
    super.connect();
    if (this.data.has('init')) return;

    const podcast = getLocalStorage(
      'podcasts',
      (podcasts) => podcasts[this.numberTarget.innerText]
    );
    if (podcast) {
      this.renderProgress(podcast);
    }

    this.data.set('init', '1');
  }

  renderProgress(podcast) {
    this.progressTarget.style.display = 'block';
    this.durationTarget.innerText = composeTime(podcast.duration);
    this.barTarget.style.width = `${(podcast.currentTime / podcast.duration) * 100}%`;
  }
}
