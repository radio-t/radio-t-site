import debounce from 'lodash/debounce';
import capitalize from 'lodash/capitalize';
import Controller from '../base_controller';
import { composeTime, getLocalStorage, parseTime, updateLocalStorage } from '../utils';
import { Events } from '../events';

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
    'volumeLevel',
    'mute',
    'unmute'
  ];

  static getState() {
    return this.state;
  }

  updateState(state) {
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

  // https://developer.mozilla.org/en-US/docs/Web/Guide/Events/Media_events
  addEventListeners() {
    ['timeupdate', 'durationchange', 'play', 'pause', 'ended']
      .forEach(event => {
        const handlerName = `on${capitalize(event)}`;
        if (this[handlerName]) this.audioTarget.addEventListener(event, this[handlerName].bind(this));
      });

    const updateLoadingState = debounce((isLoading) => this.element.classList.toggle('player-loading', isLoading), 500);
    const eventsLoadingOn = ['seeking', 'waiting', 'loadstart'];
    const eventsLoadingOff = ['playing', 'seeked', 'canplay', 'loadeddata', 'error'];
    eventsLoadingOn.forEach(event => this.audioTarget.addEventListener(event, updateLoadingState.bind(this, true)));
    eventsLoadingOff.forEach(event => this.audioTarget.addEventListener(event, updateLoadingState.bind(this, false)));

    const debugEvents = [
      "abort",
      "canplay",
      "canplaythrough",
      "durationchange",
      "emptied",
      "encrypted",
      "ended",
      "error",
      "interruptbegin",
      "interruptend",
      "loadeddata",
      "loadedmetadata",
      "loadstart",
      "mozaudioavailable",
      "pause",
      "play",
      "playing",
      "progress",
      "ratechange",
      "seeked",
      "seeking",
      "stalled",
      "suspend",
      "timeupdate",
      "volumechange",
      "waiting"
    ];
    debugEvents.forEach(event => this.audioTarget.addEventListener(event, (e) => this.debug('audio event', event, e)));

    window.addEventListener('beforeunload', (e) => {
      const isPlaying = this.constructor.state.src && !this.constructor.state.paused;
      if (!isPlaying) return;
      window.localStorage.setItem('player', JSON.stringify({
        ...this.detail,
        currentTime: this.audioTarget.currentTime,
        duration: this.audioTarget.duration,
      }));
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
    this.detail = detail;
    if (this.audioTarget.src !== detail.src) {
      this.resetUI();
      this.element.classList.remove('d-none');
      this.audioTarget.src = detail.src;
      this.updateState({src: detail.src});
      this.linkTargets.forEach((link) => link.href = detail.url);
      this.coverTarget.style.backgroundImage = detail.image;
      this.coverTarget.classList.toggle('cover-image-online', !!detail.online);
      this.numberTarget.textContent = detail.number;
      this.audioTarget.load();
      this.audioTarget.volume = 1;
      this.storedVolumeLevel = 100;
      if (!detail.online) {
        if (!detail.timeLabel) {
          const podcast = getLocalStorage('podcasts', podcasts => podcasts[this.numberTarget.innerText]);
          if (podcast) {
            detail.timeLabel = composeTime(podcast.currentTime);
          }
        }
      }
      this.once(this.audioTarget, 'canplay', () => this.setTimeLabel(detail.timeLabel));
      return true;
    }
    return false;
  }

  resetUI() {
    this.showPlayButton();
    this.updateCurrentTime(0, true);
    this.durationTarget.textContent = '--:--:--';
  }

  setTimeLabel(timeLabel) {
    if (timeLabel) {
      this.updateCurrentTime(parseTime(timeLabel));
    }
    return !!timeLabel;
  }

  updateCurrentTime(time, onlyUI = false) {
    if (!onlyUI) this.audioTarget.currentTime = time;
    this.seekTarget.value = time;
    this.currentTimeTarget.textContent = composeTime(time);
  }

  playPause() {
    if (this.audioTarget.paused) {
      return this.audioTarget.play();
    } else {
      return this.audioTarget.pause();
    }
  }

  seekBackward() {
    this.updateCurrentTime(this.audioTarget.currentTime - 15);
  }

  seekForward() {
    this.updateCurrentTime(this.audioTarget.currentTime + 15);
  }

  seeking(e) {
    this.isSeeking = true;
    this.updateCurrentTime(e.target.value, true);
    this.seekLaterTimeout = this.seekLater(e);
  }

  seekLater = debounce(this.seek, 100);

  seek(e) {
    this.isSeeking = false;
    if (this.audioTarget.duration) this.updateCurrentTime(e.target.value);
    clearTimeout(this.seekLaterTimeout);
  }

  volumeChanging(e) {
    this.setVolume(Number.parseInt(e.target.value, 10));
  }
  volumeChange(e) {
    const volume = Number.parseInt(e.target.value, 10);
    this.setVolume(volume);
    if (volume !== 0) {
      this.storedVolumeLevel = volume;
    }
  }
  setVolume(volumeLevel) {
    this.isMuted = volumeLevel === 0;
    this.muteTarget.classList.toggle('d-none', this.isMuted);
    this.unmuteTarget.classList.toggle('d-none', !this.isMuted);
    this.audioTarget.volume = volumeLevel / 100;
  }

  close() {
    window.localStorage.removeItem('player');
    this.element.classList.add('d-none');
    this.audioTarget.src = '';
    this.updateState({src: null, paused: null});
    this.detail = {};
  }

  onTimeupdate() {
    if (this.isSeeking) return;
    this.seekTarget.value = this.audioTarget.currentTime;
    this.currentTimeTarget.textContent = composeTime(this.audioTarget.currentTime);

    if (!this.detail.online) {
      if (this.detail.number) {
        const podcastData = {
          currentTime: this.audioTarget.currentTime,
          duration: this.audioTarget.duration,
        };
        updateLocalStorage('podcasts', (podcasts) => {
          podcasts[this.detail.number] = podcastData;
          return podcasts;
        });
        Events.publish(`playing-progress-${this.detail.number}`, podcastData);
      }
    }
  }

  onDurationchange() {
    this.seekTarget.max = this.audioTarget.duration;
    this.durationTarget.textContent = composeTime(this.audioTarget.duration);
  }

  onPlay() {
    this.showPlayButton(false);
    this.updateState({paused: false});
  }

  onPause() {
    this.showPlayButton(true);
    this.updateState({paused: true});
  }

  showPlayButton(paused = false) {
    this.playTarget.classList.toggle('d-none', !paused);
    this.pauseTarget.classList.toggle('d-none', paused);
  }

  onEnded() {
    // @todo:
  }

  muteUnmute() {
    if (this.isMuted) {
      this.setVolume(this.storedVolumeLevel);
      this.volumeLevelTarget.value = this.storedVolumeLevel;
    } else {
      this.setVolume(0);
      this.volumeLevelTarget.value = 0;
    }
  }
}
