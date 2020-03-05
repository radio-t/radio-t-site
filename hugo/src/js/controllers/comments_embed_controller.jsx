import React from 'react';
import { render, unmountComponentAtNode } from 'react-dom';
import Controller from '../base_controller';
import Remark from '../components/remark';
import Turbolinks from 'turbolinks';

export default class extends Controller {
  initialize() {
    super.initialize();
    this.render = this.render.bind(this);
  }

  connect() {
    super.connect();
    window.remark_config = window.remark_config || {};
    window.remark_config.url = 'https://radio-t.com' + location.pathname;

    this.render();
    document.addEventListener('theme:change', this.render);
  }

  render() {
    const theme = window.RADIOT_THEME === 'dark' ? 'dark' : 'light';

    render((<Remark
      site_id={window.remark_config.site_id}
      url={'https://radio-t.com' + location.pathname}
      page_title={window.remark_config.page_title}
      theme={theme}
      locale={window.remark_config.locale}
    />), this.element);
  }

  disconnect() {
    super.disconnect();
    document.removeEventListener('theme:change', this.render);
    unmountComponentAtNode(this.element);
  }
}
