import React, { Component } from 'react';
import { distanceInWordsStrict, format, parse } from 'date-fns';
import locale from 'date-fns/locale/ru';
import Visibility from 'visibilityjs';
import http from 'axios';

const COMMENT_NODE_CLASSNAME_PREFIX = 'remark42__comment-';

function getTextSnippet(html) {
  const LENGTH = 120;
  const tmp = document.createElement('div');
  tmp.innerHTML = html.replace('</p><p>', ' ');

  const result = tmp.innerText || '';
  const snippet = result.substr(0, LENGTH);

  return snippet.length === LENGTH && result.length !== LENGTH ? `${snippet}...` : snippet;
}

const Comment = function ({comment}) {
  const now = new Date();
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
        <h5 className="m-0 small font-weight-bold">
          {comment.user.name}
          {comment.user.verified && <div className="last-comments-comment__verification"></div>}
        </h5>
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

let comments;

class LastComments extends Component {
  constructor() {
    super();
    this.state = {comments: []};
  }

  componentDidMount() {
    const min = 60 * 1000;
    this.visibilityInterval = Visibility.every(min / 2, 5 * min, () => {
      this.updateComments();
    });

    this.updateComments();
  }

  async updateComments() {
    const {data} = await http.get('https://remark42.radio-t.com/api/v1/last/30', {params: {site: 'radiot'}});
    this.setState({comments: data});
  }

  componentWillUnmount() {
    Visibility.stop(this.visibilityInterval);
  }

  render() {
    return <div className="last-comments-list">{this.state.comments.map((comment) =>
      <Comment comment={comment} key={comment.id}/>,
    )}</div>;
  }
}

export default LastComments;
