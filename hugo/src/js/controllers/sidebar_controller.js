import Controller from '../base_controller';

export default class extends Controller {
  static targets = ['scroll'];

  toggle() {
    this.element.classList.toggle('sidebar-open');
  }

  // data-action="scroll->sidebar#preventScrollPropagation wheel->sidebar#preventScrollPropagation"
  preventScrollPropagation(e) {
    // const n = e.deltaY || e.detail || e.wheelDelta;
    // const s = this.scrollTarget;
    // if (n < 0 && 0 === s.scrollTop || (n > 0 && s.scrollTop >= s.scrollHeight - s.clientHeight)) {
    //   e.preventDefault();
    //   e.stopPropagation();
    // }
  }
}
