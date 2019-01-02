import { Controller } from 'stimulus';
import padStart from 'lodash/padStart';
// import {Howl, Howler} from 'howler';

/**
 * @property {HTMLAudioElement} audioTarget
 */
export default class extends Controller {
  static targets = ['audio', 'seek', 'play', 'pause', 'currentTime'];

  initialize() {
    // this.audioTarget.on
  }

  async timeLabel(src, time) {
    this.audioTarget.src = src;
    this.audioTarget.currentTime = this.parseTime(time);

    return await this.audioTarget.play();
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
    while (pieces.length < 3) pieces.push(0)
    return pieces.reverse().map((t) => padStart(t, 2, '0')).join(':');
  }

  playPause(src) {
    if (src) {
      if (this.audioTarget.src === src) {
        return this.playPause();
      }
      this.audioTarget.src = src;
      this.data.set('src', src);
      return this.audioTarget.play();
    } else {
      if (this.audioTarget.paused) {
        return this.audioTarget.play();
      } else {
        return this.audioTarget.pause();
      }
    }
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

  timeupdate() {
    this.seekTarget.style.left = `${this.audioTarget.currentTime / this.audioTarget.duration * 100}%`;
    this.currentTimeTarget.textContent = this.composeTime(this.audioTarget.currentTime);
  }
}
