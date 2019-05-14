import React from 'react';
import { render } from 'react-dom';
import Controller from '../base_controller';
import LastComments from '../components/LastComments';

export default class extends Controller {
  async initialize() {
    super.initialize();

    render((<LastComments/>), this.element);
  }
}
