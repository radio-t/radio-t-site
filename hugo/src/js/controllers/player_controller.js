import { Controller } from 'stimulus';
import padStart from 'lodash/padStart';

function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

/**
 * @property {Audio} audioTarget
 */
export default class extends Controller {
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

  initialize() {
    this.addEventListeners();
  }

  addEventListeners() {
    const events = ['timeupdate', 'durationchange', 'play', 'pause', 'ended'];
    events.forEach((event) => {
      const handlerName = `on${capitalizeFirstLetter(event)}`;
      if (this[handlerName]) this.audioTarget.addEventListener(event, this[handlerName].bind(this));
    });
  }

  onTimeupdate() {
    this.seekTarget.style.left = `${this.audioTarget.currentTime / this.audioTarget.duration * 100}%`;
    this.currentTimeTarget.textContent = this.composeTime(this.audioTarget.currentTime);
  }

  onDurationchange() {
    this.durationTarget.textContent = this.composeTime(this.audioTarget.duration);
  }

  onPlay() {
    this.playTarget.classList.add('d-none');
    this.pauseTarget.classList.remove('d-none');
  }

  onPause() {
    this.playTarget.classList.remove('d-none');
    this.pauseTarget.classList.add('d-none');
  }

  onEnded() {
    // @todo:
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

  seek(e) {
    if (this.audioTarget.duration) {
      this.audioTarget.currentTime = e.layerX / e.target.clientWidth * this.audioTarget.duration;
    }
  }

}
