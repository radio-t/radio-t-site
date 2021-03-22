import { h } from 'preact';
import { useCallback, useEffect, useState } from 'preact/hooks';
import http from 'axios';
import locale from 'date-fns/locale/ru';
import formatDistanceToNowStrict from 'date-fns/formatDistanceToNow';
import parseISO from 'date-fns/parseISO'
import { getTextSnippet } from '../utils';
const COMMENT_NODE_CLASSNAME_PREFIX = 'remark42__comment-';

function Comment({ comment }) {
  const date = parseISO(comment.time);
  const href = `${new URL(comment.locator.url).pathname}#${COMMENT_NODE_CLASSNAME_PREFIX}${
    comment.id
  }`;

  return (
    <div className="last-comments-list-item">
      <div className="mb-3 text-content">
        <a href={href}>{comment.title}</a>
        <small className="text-muted"> &rarr;</small>
      </div>
      <div className="mb-2 media last-comments-header align-items-center">
        <div className="last-comments-avatar mr-2">
          {comment.user.picture && (
            <img
              className="last-comments-avatar-image"
              src={comment.user.picture}
              alt={comment.user.name}
              loading="lazy"
              width={28}
              height={28}
            />
          )}
        </div>
        <div className="media-body">
          <h5 className="m-0 small font-weight-bold">
            {comment.user.name}
            {comment.user.verified && <div className="last-comments-comment__verification" />}
          </h5>
          <a
            href={href}
            title={format(date, 'dd MMM yyyy, HH:mm', { locale })}
            className="small text-muted"
          >
            {formatDistanceToNowStrict(date, { locale, addSuffix: true })}
          </a>
        </div>
      </div>
      <div className="last-comments-text">
        <a href={href}>{getTextSnippet(comment.text)}</a>
      </div>
    </div>
  );
}

function LastComments() {
  const [comments, setComments] = useState([]);

  const updateComments = useCallback(async () => {
    try {
      const { data } = await http.get('https://remark42.radio-t.com/api/v1/last/30', {
        params: { site: 'radiot' },
      });
      setComments(data);
    } catch (e) {
      //
    }
  }, [setComments]);

  useEffect(() => {
    updateComments();
    document.addEventListener('turbolinks:visit', updateComments);
    return () => document.removeEventListener('turbolinks:visit', updateComments);
  }, [updateComments]);

  return (
    <div className="last-comments-list">
      {comments.map((comment) => (
        <Comment comment={comment} key={comment.id} />
      ))}
    </div>
  );
}

export default LastComments;
