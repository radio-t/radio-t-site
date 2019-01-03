import { Controller } from 'stimulus';
import Player from './player_controller';

/**
 * @property playButtonTarget
 * @property {Audio} audioTarget
 */
export default class extends Controller {
  static targets = ['playButton', 'number', 'cover', 'audio'];

  initialize() {
    // Set up audio target
    const audio = this.element.querySelector('audio');
    if (audio && audio.src) {
      audio.dataset.target = `${this.identifier}.audio`;
      this.playButtonTarget.classList.remove('d-none');
      this.element.classList.add('has-audio');
    }
  }

  // connect() {
  //   this.element.dispatchEvent(new CustomEvent('podcast-connected', {bubbles: true}));
  // }

  play(e, timeLabel = null) {
    e.preventDefault();
    e.stopPropagation();

    // see https://github.com/stimulusjs/stimulus/issues/200#issuecomment-434731830
    this.element.dispatchEvent(new CustomEvent('podcast-play', {
      bubbles: true,
      detail: {
        ...this.getPodcastInfo(),
        ...(timeLabel ? {timeLabel} : {}),
      },
    }));
  }

  goToTimeLabel(e) {
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
}
