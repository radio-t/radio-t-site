import Visibility from 'visibilityjs';
import Player from './player_controller';
import Controller from '../base_controller';
import { timer } from '../timer';

// const STREAM_SRC = 'http://rfcmedia.streamguys1.com/MusicPulse.mp3'; // demo
const STREAM_SRC = 'https://stream.radio-t.com/';

/**
 * This handles online page and online banner.
 *
 * @property {HTMLElement} labelTarget
 * @property {HTMLElement} timeTarget
 */
export default class extends Controller {
  static targets = ['label', 'time'];

  connect() {
    super.connect();

    this.setupTimer();
    this.fetchPlayingState();
  }

  disconnect() {
    super.disconnect();
    Visibility.stop(this.visibilityInterval);
  }

  setupTimer() {
    const tick = () => {
      const t = timer();
      this.timeTarget.innerHTML = t.html;
      this.element.classList.toggle('is-online', t.isOnline);
    };

    tick();
    this.visibilityInterval = Visibility.every(1000, 60000, tick);
  }

  fetchPlayingState() {
    this.element.classList.toggle('playing', this.isCurrentlyPlaying());
  }

  isCurrentlyPlaying() {
    return (
      Player.getState().src === this.getPodcastInfo().src
      && Player.getState().paused === false
    );
  }

  play(e) {
    e.preventDefault();
    e.stopPropagation();

    this.dispatchEvent(this.element, new CustomEvent('podcast-play', {
      bubbles: true,
      detail: this.getPodcastInfo(),
    }));
  }

  getPodcastInfo() {
    return {
      src: STREAM_SRC,
      url: '/online',
      image: null,
      number: 'Online',
      online: true,
    };
  }
}
