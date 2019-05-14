import Controller from '../base_controller';

export default class extends Controller {
  toggle() {
    this.element.classList.toggle('sidebar-open');
  }
}
