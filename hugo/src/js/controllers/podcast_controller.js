import Controller from '../base_controller';
import Player from './player_controller';
import { addTimeToURL } from '../utils'

/**
 * @property playButtonTarget
 * @property {Audio} audioTarget
 */
export default class extends Controller {
  static targets = ['playButton', 'number', 'cover', 'audio'];

  initialize() {
    super.initialize();

    // Set up audio target
    const audio = this.element.querySelector('audio');

    if (audio && audio.src) {
      audio.dataset.target = `${this.identifier}.audio`;
      this.playButtonTarget.classList.remove('d-none');
      this.element.classList.add('has-audio');
    }

    this.fetchPlayingState();
  }

  fetchPlayingState() {
    this.element.classList.toggle('playing', this.isCurrentlyPlaying());
  }

  isCurrentlyPlaying() {
    return (
      Player.getState().src === this.getPodcastInfo().src && Player.getState().paused === false
    );
  }

  play(e, timeLabel = null) {
    e.preventDefault();
    e.stopPropagation();

    this.dispatchEvent(
      this.element,
      new CustomEvent('podcast-play', {
        bubbles: true,
        detail: {
          ...this.getPodcastInfo(),
          ...(timeLabel ? { timeLabel } : {}),
        },
      })
    );

    setTimeout(() => this.fetchPlayingState(), 0);
  }

  goToTimeLabel(e) {
    // add each seek time to URL as t?=00:00:00
    this.podcastAddTimeToURL(e);
    this.play(e, e.target.textContent);
  }

  getPodcastInfo() {
    return {
      src: this.audioTarget.src,
      url: this.data.get('url'),
      image: this.coverTarget.style.backgroundImage,
      number: this.numberTarget.textContent,
    };
  }

  podcastAddTimeToURL(e) {
    let podcastPathname = this.data.get('url');
    addTimeToURL(podcastPathname, e.target.textContent);
  }
}
