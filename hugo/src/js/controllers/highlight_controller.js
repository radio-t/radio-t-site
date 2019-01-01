import { Controller } from 'stimulus';
import hljs from '../highlight'

export default class extends Controller {
  initialize() {
    this.element.querySelectorAll('pre code').forEach((code) => hljs.highlightBlock(code))
  }
}
