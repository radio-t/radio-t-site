import Controller from '../base_controller';
import likely from 'ilyabirman-likely';

export default class extends Controller {
  connect() {
    super.connect();
    likely.initiate(this.element, {
      counters: false,
    });
  }
}
