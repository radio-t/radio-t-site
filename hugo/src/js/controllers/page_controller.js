import fastclick from 'fastclick';
import Controller from '../base_controller';
import Player from './player_controller';

export default class extends Controller {
  static targets = [
    'player',
    'podcast',
  ];

  initialize() {
    super.initialize();
    fastclick.attach(this.element);
  }

  updatePodcasts(e) {
    this.podcastTargets.forEach((podcast) => {
      this.dispatchEvent(podcast, new CustomEvent('player-state', {bubbles: false}))
    });
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
