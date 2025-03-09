import map from 'lodash/map';
import uniq from 'lodash/uniq';
import filter from 'lodash/filter';
import lozad from 'lozad';
import imagesLoaded from 'imagesloaded';
import Controller from '../base_controller';
import { fetchJSON } from '../http-get';

const limit = 30;

export default class extends Controller {
  connect() {
    super.connect();
    lozad(this.element, {
      load: async () => {
        this.element.classList.add('loaded');
        const data = await this.getComments();
        const pictures = uniq(filter(map(data.comments, 'user.picture'))).reverse();
        this.element.innerHTML = '';
        pictures.slice(0, limit).forEach((picture, index) => {
          if (!picture) return;
          const div = document.createElement('DIV');
          div.style.backgroundImage = `url('${picture}')`;
          div.classList.add('comments-counter-avatars-item');
          div.style.transitionDelay = `${(limit - index) * 20}ms`;
          this.element.appendChild(div);
        });
        // todo limit waiting time before showing images
        imagesLoaded(this.element, { background: '.comments-counter-avatars-item' }, () => {
          this.reflow();
          this.element.classList.remove('loaded');
        });
      },
    }).observe();
  }

  async getComments() {
    const url = this.data.get('url');

    return fetchJSON(
      `https://remark42.radio-t.com/api/v1/find?url=https://radio-t.com${url}&sort=-time&site=radiot&view=user`
    );
  }
}
