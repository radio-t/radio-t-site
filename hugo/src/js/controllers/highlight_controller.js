import Controller from '../base_controller';
import hljs from '../highlight'

export default class extends Controller {
  connect() {
    super.connect();
    this.element.querySelectorAll('pre code').forEach((code) => hljs.highlightBlock(code))
  }
}
