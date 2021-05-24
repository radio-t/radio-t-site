import Controller from '../base_controller';
import { composeTime, getLocalStorage } from '../utils';

export default class extends Controller {
  static targets = ['bar', 'number', 'duration', 'progress'];

  get podcastNumber() {
    return this.numberTarget.dataset.number;
  }

  initialize() {
    super.initialize();
    this.subscribe(`playing-progress-${this.podcastNumber}`, (podcast) => {
      this.renderProgress(podcast);
    });
  }

  connect() {
    super.connect();
    if (this.data.has('init')) return;

    const podcast = getLocalStorage('podcasts', (podcasts) => podcasts[this.podcastNumber]);
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
