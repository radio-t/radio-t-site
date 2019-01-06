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
}
