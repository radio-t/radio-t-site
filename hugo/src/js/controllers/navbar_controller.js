import Controller from '../base_controller';

export default class extends Controller {
  static targets = ['backdrop', 'scroll'];

  initialize() {
    super.initialize();
    document.addEventListener('turbolinks:before-cache', () => {
      this.element.classList.remove('navbar-open');
    });
    this.element.querySelectorAll('.nav-item').forEach((item, index) => {
      item.style.transitionDuration = `${index * 30 + 150}ms`;
    });
  }

  connect() {
    super.connect();
    this.element.classList.remove('navbar-open');

    // Mark active nav item
    for (let item of this.element.querySelectorAll('.nav-link')) {
      if (item.hasAttribute('no-active')) continue;
      item.classList.toggle('active', window.location.href.startsWith(item.href));
    }
  }

  toggle() {
    if (!this.element.classList.contains('navbar-open')) {
      this.scrollTarget.scrollTo(0, 0);
    }
    this.element.classList.toggle('navbar-open');
  }

  closeFromBackdrop(e) {
    if (e.target === this.backdropTarget) {
      this.element.classList.remove('navbar-open');
    }
  }
}
