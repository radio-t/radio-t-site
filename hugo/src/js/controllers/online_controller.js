import Visibility from 'visibilityjs';
import Controller from '../base_controller';
import Player from './player_controller';

// const STREAM_SRC = 'http://rfcmedia.streamguys1.com/MusicPulse.mp3'; // demo
const STREAM_SRC = 'http://stream.radio-t.com/';
const showTime = {
  day: 6,
  hours: 23,
  minutes: 0,
};

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
      const timer = this.timer();
      this.timeTarget.innerHTML = timer.html;
      this.element.classList.toggle('is-online', timer.isOnline);
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

  timer() {
    function getUnits(value, units) {
      return (/^[0,2-9]?[1]$/.test(value)) ? units[0] : ((/^[0,2-9]?[2-4]$/.test(value)) ? units[1] : units[2]);
    }

    function padTime(n) {
      return ('0' + n).slice(-2);
    }

    const timeInMoscow = new Date();
    timeInMoscow.setMinutes(timeInMoscow.getMinutes() + timeInMoscow.getTimezoneOffset() + 3 * 60);

    const nextShow = new Date(timeInMoscow);
    nextShow.setDate(nextShow.getDate() + showTime.day - nextShow.getDay());
    nextShow.setHours(showTime.hours, showTime.minutes, 0, 0);

    const totalSeconds = Math.floor((nextShow - timeInMoscow) / 1000);

    const isOnline = totalSeconds <= 0;
    if (isOnline) {
      return {isOnline, html: 'Мы в эфире!'};
    }

    let seconds = totalSeconds % 60,
      minutes = Math.round((totalSeconds - seconds) / 60) % 60,
      hours = Math.round((totalSeconds - seconds - minutes * 60) / 3600),
      days = (hours - hours % 24) / 24;

    hours %= 24;

    let html = '';
    const daysList = ['день', 'дня', 'дней'];

    if (days > 0) {
      html += days + ' ' + getUnits(days, daysList) + ' ';
    }

    html += `${padTime(hours)}:${padTime(minutes)}<span style="opacity: .5;">:${padTime(seconds)}</span>`;

    return {isOnline, html};
  }
}
