import $script from 'scriptjs';
import Controller from '../base_controller';

export default class extends Controller {
  initialize() {
    super.initialize();
    ((window.gitter = {}).chat = {}).options = {
      room: 'radio-t/chat',
      // activationElement: this.element,
      // targetElement: '.gitter-sidecar-container',
    };
    $script.get('https://sidecar.gitter.im/dist/sidecar.v1.js', () => {});
  }
}
