import http from 'axios';
import Visibility from 'visibilityjs';
import locale from 'date-fns/locale/ru';
import { distanceInWordsStrict, format, parse } from 'date-fns';
import React, { useCallback, useEffect, useState } from 'react';
import { getTextSnippet } from '../utils';

const COMMENT_NODE_CLASSNAME_PREFIX = 'remark42__comment-';

function Comment({comment}) {
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
}

function LastComments() {
  const [comments, setComments] = useState([]);

  const updateComments = useCallback(async () => {
    try {
      const {data} = await http.get('https://remark42.radio-t.com/api/v1/last/30', {params: {site: 'radiot'}});
      setComments(data);
    } catch (e) {
      //
    }
  }, [setComments]);

  useEffect(() => {
    updateComments();
    if (process.env.NODE_ENV !== 'development') {
      const visibilityInterval = Visibility.every(60 * 1000, updateComments);
      return () => Visibility.stop(visibilityInterval);
    }
  }, [updateComments]);

  return <div className="last-comments-list">{comments.map((comment) =>
    <Comment comment={comment} key={comment.id}/>,
  )}</div>;
}

export default LastComments;
