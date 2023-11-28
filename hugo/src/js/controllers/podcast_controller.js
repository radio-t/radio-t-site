import Controller from '../base_controller';
import Player from './player_controller';

/**
 * @property playButtonTarget
 * @property {Audio} audioTarget
 */
export default class extends Controller {
  static targets = ['playButton', 'number', 'cover', 'audio', 'icon'];

  get isCurrentlyPlaying() {
    const state = Player.getState();

    return state.src === this.audioTarget.src && state.paused === false;
  }

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
    this.playButtonTarget.classList.toggle('podcast-cover_playing', this.isCurrentlyPlaying);
    this.iconTarget.classList.toggle('podcast-cover-icon_play', !this.isCurrentlyPlaying);
    this.iconTarget.classList.toggle('podcast-cover-icon_pause', this.isCurrentlyPlaying);
    this.iconTarget.firstChild.setAttribute(
      'xlink:href',
      this.isCurrentlyPlaying ? '#icon-pause' : '#icon-play'
    );
  }

  toggle(e, timeLabel = null) {
    e.preventDefault();
    e.stopPropagation();
    const detail = {
      src: this.audioTarget.src,
      url: this.data.get('url'),
      image: this.coverTarget.currentSrc,
      number: this.playButtonTarget.dataset.number,
      timeLabel,
    };
    console.log(detail.image);
    this.dispatchEvent(this.element, new CustomEvent('podcast-play', { bubbles: true, detail }));

    setTimeout(() => this.fetchPlayingState(), 0);
  }

  goToTimeLabel(e) {
    this.play(e, e.target.textContent);
  }
}
