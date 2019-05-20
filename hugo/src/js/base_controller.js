import { Controller } from 'stimulus';
import { Events } from './events';
import debug from './debug';

export default class extends Controller {
  subscriptions = {};
  unsubscribe = [];

  subscribe(eventName, handler) {
    this.subscriptions[eventName] = this.subscriptions[eventName] || [];
    this.subscriptions[eventName].push(handler);
  }

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
    for (const eventName in this.subscriptions) {
      for (const handler of this.subscriptions[eventName]) {
        const {unsubscribe} = Events.subscribe(eventName, handler);
        this.unsubscribe.push(unsubscribe);
      }
    }
  }

  disconnect() {
    this.debug('disconnect');
    super.disconnect();
    for (const unsubscribe of this.unsubscribe) {
      unsubscribe();
    }
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
