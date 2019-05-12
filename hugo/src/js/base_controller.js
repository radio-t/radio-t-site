import { Controller } from 'stimulus';
import debug from './debug';

export default class extends Controller {
  debug(...args) {
    debug(`stimulus:${this.identifier}`, ...[...args, this.element]);
  }

  constructor(...args) {
    const instance = super(...args);
    this.debug('constructor');
    return instance;
  }

  initialize() {
    this.debug('initialize');
    super.initialize();
  }

  connect() {
    this.debug('connect');
    super.connect();
  }

  disconnect() {
    this.debug('disconnect');
    super.disconnect();
  }

  // see https://github.com/stimulusjs/stimulus/issues/200#issuecomment-434731830
  dispatchEvent(element, event) {
    this.debug('dispatch event', event, element);
    element.dispatchEvent(event);
  }

  once(el, event, handler) {
    function realHandler(...args) {
      el.removeEventListener(event, realHandler);
      handler(...args);
    }

    el.addEventListener(event, realHandler);
  }

  reflow(element) {
    element = element || this.element;
    return element.offsetHeight;
  }
}
