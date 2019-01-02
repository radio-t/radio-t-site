import { Controller } from 'stimulus';
import Player from './player_controller';

export default class extends Controller {
  static targets = [
    'player',
    'podcast',
  ];

  initialize() {
    // console.log('initialize');
    // this.audioTarget.on
    // this.updatePodcasts();
  }

  updatePodcasts(e) {
    // console.log(e);
    // const targets = e ? [e.target] : this.podcastTargets;
    // targets.forEach((podcast) => {
    //   console.log(podcast);
    // });
  }

  playPodcast({detail}) {
    this.getPlayerController().playPodcast(detail);
  }

  /**
   * @return {Player}
   */
  getPlayerController() {
    return this.application.getControllerForElementAndIdentifier(this.playerTarget, 'player');
  }
}
