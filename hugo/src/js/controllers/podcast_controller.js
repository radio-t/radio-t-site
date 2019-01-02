import { Controller } from 'stimulus';
import filter from 'lodash/filter';
import find from 'lodash/find';

/**
 * @property playButtonTarget
 * @property audioTarget
 */
export default class extends Controller {
  static targets = ['audio', 'playButton'];

  audioSrc;

  initialize() {
    const audio = this.element.querySelector('audio');
    if (audio) {
      audio.dataset.target = `${this.identifier}.audio`;
      this.element.classList.add('has-audio');
    }
    this.playButtonTarget.classList.remove('d-none');
  }

  jumpTime(e) {
    console.log(`Jump to time ${e.target.textContent}`);
  }

  play(e) {
    console.log('play');
    e.preventDefault();
    e.stopPropagation();
  }
}
