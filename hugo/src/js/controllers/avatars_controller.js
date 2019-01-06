import Controller from '../base_controller';
import axios from 'axios';
import { cacheAdapterEnhancer, throttleAdapterEnhancer } from 'axios-extensions';
import map from 'lodash/map';
import uniq from 'lodash/uniq';
import filter from 'lodash/filter';

const http = axios.create({
  adapter: throttleAdapterEnhancer(cacheAdapterEnhancer(axios.defaults.adapter)),
});

export default class extends Controller {
  async initialize() {
    super.initialize();
    if (this.data.get('initialized')) return;
    this.data.set('initialized', '1');
    const {data} = await this.getComments();
    const pictures = uniq(filter(map(data.comments, 'user.picture')));
    pictures.slice(0, 10).forEach((picture) => {
      if (!picture) return;
      const div = document.createElement('DIV');
      div.style.backgroundImage = `url('${picture}')`;
      div.classList.add('comments-counter-avatars-item');
      this.element.append(div);
    })
  }

  getComments() {
    return http.get('https://remark42.radio-t.com/api/v1/find?url=https://radio-t.com/p/2018/12/29/podcast-630/&sort=-time&site=radiot');
  }
}
