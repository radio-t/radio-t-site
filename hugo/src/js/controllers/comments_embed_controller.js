import { h, render } from 'preact';
import Controller from '../base_controller';
import Remark from '../components/remark';

export default class extends Controller {
  connect() {
    super.connect();
    window.remark_config = window.remark_config || {};
    window.remark_config.url = 'https://radio-t.com' + location.pathname;

    this.root = render((<Remark
      site_id={window.remark_config.site_id}
      url={'https://radio-t.com' + location.pathname}
    />), this.element);
  }
  disconnect() {
    super.disconnect();
    render(null, this.element);
  }
}
