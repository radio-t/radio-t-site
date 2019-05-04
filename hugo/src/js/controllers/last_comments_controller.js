import { h, render } from 'preact';
import { format, parse, distanceInWordsToNow, distanceInWordsStrict } from 'date-fns';
import locale from 'date-fns/locale/ru';
import Visibility from 'visibilityjs';
import Controller from '../base_controller';
import http from 'axios';
// import http from '../http-client';

const COMMENT_NODE_CLASSNAME_PREFIX = 'remark42__comment-';

export default class extends Controller {
  async initialize() {
    super.initialize();

    await this.updateComments();

    const min = 60 * 1000;
    Visibility.every(min / 2, 10 * min, async () => {
      await this.updateComments();
    });
  }

  async updateComments() {
    // Set min height to preserve scroll position
    const height = this.element.offsetHeight;
    const prevMinHeight = this.element.style.minHeight;
    this.element.style.minHeight = `${height}px`;

    const {data} = await http.get('https://remark42.radio-t.com/api/v1/last/30', {params: {site: 'radiot'}});
    this.element.innerHTML = '';
    render((<LastComments comments={data}/>), this.element);

    this.element.style.minHeight = prevMinHeight;
  }
}

function getTextSnippet(html) {
  const LENGTH = 100;
  const tmp = document.createElement('div');
  tmp.innerHTML = html.replace('</p><p>', ' ');

  const result = tmp.innerText || '';
  const snippet = result.substr(0, LENGTH);

  return snippet.length === LENGTH && result.length !== LENGTH ? `${snippet}...` : snippet;
}

const LastComments = function ({comments}) {
  const now = new Date();

  const Comment = function ({comment}) {
    const date = parse(comment.time);
    const avatarStyle = comment.user.picture ? {backgroundImage: `url('${comment.user.picture}')`} : {};
    const href = (new URL(comment.locator.url)).pathname + `#${COMMENT_NODE_CLASSNAME_PREFIX}${comment.id}`;

    return <div className="last-comments-list-item">
      <div className="mb-3 text-content">
        <a href={href}>{comment.title}</a>
        <small className="text-muted"> &rarr;</small>
      </div>
      <div className="mb-2 media last-comments-header align-items-center">
        <div className="last-comments-avatar mr-2">
          <div className="last-comments-avatar-image" style={avatarStyle}/>
        </div>
        <div className="media-body">
          <h5 className="m-0 small font-weight-bold">{comment.user.name}</h5>
          <a
            href={href}
            title={format(date, 'DD MMM YYYY, HH:mm', {locale})}
            className="small text-muted"
          >
            {distanceInWordsStrict(now, date, {locale, addSuffix: true})}
          </a>
        </div>
      </div>
      <div className="last-comments-text">
        <a href={href}>{getTextSnippet(comment.text)}</a>
      </div>
    </div>;
  };

  return <div className="last-comments-list">{comments.map((comment) =>
    <Comment comment={comment}/>,
  )}</div>;
};
