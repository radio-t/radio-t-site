import Controller from '../base_controller';

/**
 * @property {HTMLElement} labelTarget
 * @property {HTMLElement} timeTarget
 */
export default class extends Controller {
  static targets = ['label', 'time'];

  connect() {
    super.connect();

    this.timeTarget.innerHTML = this.setShowTimer();
    window.setInterval(() => {
      this.timeTarget.innerHTML = this.setShowTimer();
    }, 1000);
  }

  play () {
    this.dispatchEvent(this.element, new CustomEvent('podcast-play', {
      bubbles: true,
      detail: {
        src: 'http://stream.radio-t.com/',
        url: '/online',
        image: null,
        number: 'Online',
        online: true,
      }
    }));
  }

  setShowTimer() {
    function getUnits(value, units) {
      return (/^[0,2-9]?[1]$/.test(value)) ? units[0] : ((/^[0,2-9]?[2-4]$/.test(value)) ? units[1] : units[2]);
    }

    function padTime(n) {
      return ('0' + n).slice(-2);
    }

    const timeInMoscow = new Date();
    // timeInMoscow.setDate(timeInMoscow.getDate() + 6 - timeInMoscow.getDay());
    // timeInMoscow.setHours(22, 30);
    timeInMoscow.setMinutes(timeInMoscow.getMinutes() + timeInMoscow.getTimezoneOffset() + 3 * 60);

    const nextShow = new Date(timeInMoscow);
    nextShow.setDate(nextShow.getDate() + 6 - nextShow.getDay());
    nextShow.setHours(23, 0, 0, 0);

    const totalSeconds = Math.floor((nextShow - timeInMoscow) / 1000);

    if (totalSeconds < 0) {
      return 'Мы в эфире!';
    }

    let seconds = totalSeconds % 60,
      minutes = Math.round((totalSeconds - seconds) / 60) % 60,
      hours = Math.round((totalSeconds - seconds - minutes * 60) / 3600),
      days = (hours - hours % 24) / 24;

    hours %= 24;

    let result = '';
    const daysList = ['день', 'дня', 'дней'];

    if (days > 0) {
      result += days + ' ' + getUnits(days, daysList) + ' ';
    }

    result += `${padTime(hours)}:${padTime(minutes)}<span style="opacity: .5;">:${padTime(seconds)}</span>`;

    return result;
  }
}
