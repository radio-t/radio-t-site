import pose from 'popmotion-pose';
import Controller from '../base_controller';
import Player from './player_controller';
import { composeTime, getLocalStorage } from '../utils';

/**
 * @property playButtonTarget
 * @property {Audio} audioTarget
 */
export default class extends Controller {
  static targets = ['playButton', 'number', 'cover', 'coverShadow', 'audio'];

  coverPoser;

  initialize() {
    super.initialize();

    // Set up audio target
    const audio = this.element.querySelector('audio');
    if (audio && audio.src) {
      audio.dataset.target = `${this.identifier}.audio`;
      this.playButtonTarget.classList.remove('d-none');
      this.element.classList.add('has-audio');
    }

    this.setupCoverAnimation();

    this.fetchPlayingState();
  }

  setupCoverAnimation() {
    const transition = {
      type: 'spring',
      stiffness: 800,
      mass: 1,
      damping: 30,
      // velocity: 1,
    };
    this.coverPoser = pose(this.coverTarget.parentElement, {
      init: {
        scale: .9,
        y: '0%',
        transition: {...transition, damping: 60},
      },
      elevated: {
        scale: 1,
        y: '-3%',
        transition,
      },
    });
    this.coverPoser.addChild(this.coverShadowTarget, {
      init: {
        y: '-10%',
        opacity: .33,
        transition: {...transition, damping: 60},
      },
      elevated: {
        y: '-3%',
        opacity: 1,
        transition,
      },
    });
  }

  fetchPlayingState() {
    this.element.classList.toggle('playing', this.isCurrentlyPlaying());
    this.coverPoser.set(this.isCurrentlyPlaying() ? 'elevated' : 'init');
  }

  isCurrentlyPlaying() {
    return (
      Player.getState().src === this.getPodcastInfo().src
      && Player.getState().paused === false
    );
  }

  play(e, timeLabel = null) {
    e.preventDefault();
    e.stopPropagation();

    this.dispatchEvent(this.element, new CustomEvent('podcast-play', {
      bubbles: true,
      detail: {
        ...this.getPodcastInfo(),
        ...(timeLabel ? {timeLabel} : {}),
      },
    }));

    setTimeout(() => this.fetchPlayingState(), 0);
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
