import { Controller } from 'stimulus';
import Player from './player_controller'
import filter from 'lodash/filter';
import find from 'lodash/find';

/**
 * @property playButtonTarget
 * @property audioTarget
 */
export default class extends Controller {
  static targets = ['playButton'];

  initialize() {
    const audio = this.element.querySelector('audio');
    if (audio) {
      audio.dataset.target = `${this.identifier}.audio`;
      this.element.classList.add('has-audio');
      this.data.set('src', audio.src);
      this.playButtonTarget.classList.remove('d-none');
    }
  }

  seek(e) {
    this.getPlayerController().timeLabel(this.data.get('src'), e.target.textContent);
    // console.log(`Jump to time ${e.target.textContent}`);
  }

  play(e) {
    e.preventDefault();
    e.stopPropagation();
    this.getPlayerController().playPause(this.data.get('src'));
  }

  /**
   * @returns {Player}
   */
  getPlayerController() {
    return this.application.getControllerForElementAndIdentifier(document.body, 'player');
  }
}
