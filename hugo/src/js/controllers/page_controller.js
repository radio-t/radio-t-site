import Controller from '../base_controller';
import {addTimeToURL} from '../utils'

export default class extends Controller {
  static state = {
    URLTime: null,
  };

  static targets = ['player', 'playerStateReceiver'];

  initialize() {
    super.initialize();
    this.getURLTimecode();
  }

  updatePodcasts() {
    this.playerStateReceiverTargets.forEach((podcast) => {
      this.dispatchEvent(podcast, new CustomEvent('player-state', { bubbles: false }));
    });
  }

  playPodcast({ detail }) {
    // set playback position from ?t=00:00:00 if present
    if (this.constructor.state.URLTime) {
      detail.timeLabel = this.constructor.state.URLTime;
      Object.assign(this.constructor.state, { "URLTime": null });
      addTimeToURL(detail.url, detail.timeLabel);
    }
    this.getPlayerController().playPodcast(detail);
  }

  /**
   * @return {Player}
   */
  getPlayerController() {
    return this.application.getControllerForElementAndIdentifier(this.playerTarget, 'player');
  }

  getURLTimecode() {
    let searchParams = new URLSearchParams(window.location.search);
    let timeFromUrl = searchParams.get("t");
    let timeRegExp = new RegExp(/[0-9]*:?[0-9]:?[0-9]+/m);
    if (timeRegExp.test(timeFromUrl)) {
      Object.assign(this.constructor.state, { "URLTime": timeFromUrl });
    }
  }
}
