import { Controller } from 'stimulus';

export default class extends Controller {
  initialize() {
    let titleElement = this.element;
    while (titleElement.children.length) titleElement = titleElement.firstElementChild;

    let title = titleElement.textContent.trim();

    // 'Темы для 630' => ['Темы для 630', 'Темы для', '630']
    let match = /(.+)\s(\d+)/gu.exec(title);
    if (match) {
      titleElement.innerHTML =
        `<span class="podcast-title-prefix">${match[1]}</span><br>
        <span class="podcast-title-number display-4" data-target="podcast.number">${match[2]}</span>`;
    } else {
      titleElement.innerHTML = `<span>${title}</span>`;
    }

    this.element.classList.add('title-processed');
    this.element.classList.remove('no-js');
  }
}
