import { Controller } from 'stimulus';
import likely from 'ilyabirman-likely';

export default class extends Controller {
  connect() {
    likely.initiate(this.element, {
      counters: false,
    });
  }
}
