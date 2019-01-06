import Controller from '../base_controller';
import padStart from 'lodash/padStart';

function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

/**
 * @property {Audio} audioTarget
 */
export default class extends Controller {
  static state = {
    src: null,
    paused: null,
  };

  static targets = [
    'audio',
    'seek',
    'play',
    'pause',
    'currentTime',
    'duration',
    'cover',
    'number',
    'link',
  ];

  static getState() {
    return this.state
  }

  updateState (state) {
    Object.assign(this.constructor.state, state);
    this.dispatchEvent(this.element, new CustomEvent('player-state', {
      detail: {state: this.constructor.state},
      bubbles: true,
    }));
  }

  initialize() {
    super.initialize();
    this.addEventListeners();
  }

  addEventListeners() {
    const events = ['timeupdate', 'durationchange', 'play', 'pause', 'ended'];
    events.forEach((event) => {
      const handlerName = `on${capitalizeFirstLetter(event)}`;
      if (this[handlerName]) this.audioTarget.addEventListener(event, this[handlerName].bind(this));
    });
  }

  playPodcast(detail) {
    if (this.loadPodcast(detail)) {
      return this.audioTarget.play();
    } else if (this.setTimeLabel(detail.timeLabel)) {
      return this.audioTarget.play();
    } else {
      return this.playPause();
    }
  }

  loadPodcast(detail) {
    if (this.audioTarget.src !== detail.src) {
      this.element.classList.remove('d-none');
      this.audioTarget.src = detail.src;
      this.updateState({src: detail.src});
      this.linkTargets.forEach((link) => link.href = detail.url);
      this.coverTarget.style.backgroundImage = detail.image;
      this.numberTarget.textContent = detail.number;
      this.setTimeLabel(detail.timeLabel);
      this.audioTarget.load();
      return true;
    }
    return false;
  }

  setTimeLabel(timeLabel) {
    if (timeLabel) {
      this.audioTarget.currentTime = this.parseTime(timeLabel);
    }
    return !!timeLabel;
  }

  playPause() {
    if (this.audioTarget.paused) {
      return this.audioTarget.play();
    } else {
      return this.audioTarget.pause();
    }
  }

  // 00:02:24 => 144
  parseTime(time) {
    return time
      .split(':')
      .reverse()
      .reduce((acc, curr, i) => acc + parseInt(curr) * Math.pow(60, i), 0);
  }

  // 144 => 00:02:24
  composeTime(time) {
    const pieces = [];
    time = parseInt(time);
    while (time) {
      pieces.push(time % 60);
      time = Math.floor(time / 60);
    }
    while (pieces.length < 3) pieces.push(0);
    return pieces.reverse().map((t) => padStart(t, 2, '0')).join(':');
  }

  seekBack() {
    this.audioTarget.currentTime -= 15;
  }

  seekForward() {
    this.audioTarget.currentTime += 15;
  }

  seeking(e) {
    this.isSeeking = true;
    this.currentTimeTarget.textContent = this.composeTime(e.target.value);
  }

  seek(e) {
    this.isSeeking = false;
    if (this.audioTarget.duration) {
      this.audioTarget.currentTime = e.target.value;
    }
  }

  close() {
    this.element.classList.add('d-none');
    this.audioTarget.src = '';
    this.updateState({src: null, paused: null});
  }

  onTimeupdate() {
    if (this.isSeeking) return;
    this.seekTarget.value = this.audioTarget.currentTime;
    this.currentTimeTarget.textContent = this.composeTime(this.audioTarget.currentTime);
  }

  onDurationchange() {
    this.seekTarget.max = this.audioTarget.duration;
    this.durationTarget.textContent = this.composeTime(this.audioTarget.duration);
  }

  onPlay() {
    this.playTarget.classList.add('d-none');
    this.pauseTarget.classList.remove('d-none');
    this.updateState({paused: false});
  }

  onPause() {
    this.playTarget.classList.remove('d-none');
    this.pauseTarget.classList.add('d-none');
    this.updateState({paused: true});
  }

  onEnded() {
    // @todo:
  }
}
