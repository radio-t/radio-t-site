import { h, render } from 'preact'
import Controller from '../base_controller';
import LastComments from '../components/LastComments';

export default class extends Controller {
  async initialize() {
    super.initialize();

    render(<LastComments/>, this.element);
  }
}
