import Controller from '../base_controller';

export default class extends Controller {
  initialize() {
    super.initialize();

    this.showSuggestTopicButton();
  }

  /**
   * first prep that goes before any podcast
   */
  showSuggestTopicButton() {
    // Only first page
    if (this.data.get('page-number') !== '1') return;

    // const children = [...this.element.children];
    for (let post of this.element.children) {
      if (post.classList.contains('posts-list-item-category-podcast')) {
        break;
      }
      if (post.classList.contains('posts-list-item-category-prep')) {
        const btn = post.querySelector('.btn-suggest-topic');
        if (btn) btn.classList.remove('d-none');
        break;
      }
    }
  }
}
