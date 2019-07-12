import Controller from '../base_controller';

import { setTheme } from "../theme-utils";

export default class extends Controller {
  toggle() {
    this.toggleTheme();
  }

  toggleTheme() {
    setTheme('dark' === window.RADIOT_THEME ? 'light' : 'dark');
  }
}
