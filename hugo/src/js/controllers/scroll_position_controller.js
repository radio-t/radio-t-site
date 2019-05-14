import Controller from '../base_controller';

/**
 * Save scroll position of element applied to
 */
export default class extends Controller {
  connect() {
    super.connect();
    this.element.scrollTop = this.data.get('scrollTop') || 0
  }

  onScroll(e) {
    this.data.set('scrollTop', e.target.scrollTop)
  }
}
