/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-prototype-builtins */
import { Component, createRef } from 'preact';

type RemarkTheme = 'light' | 'dark';

type Props = {
  host: string;
  site_id: string;
  page_title: string;
  url?: string;
  id?: string;
  className?: string;
  theme?: RemarkTheme;
  locale: string;
};

export default class Remark extends Component<Props> {
  protected ref = createRef<HTMLDivElement>();
  protected receiveMessages?: (
    event: Event & {
      data?: unknown;
    }
  ) => void;
  protected postHashToIframe?: (
    event: Event & {
      newURL: string;
    }
  ) => void;
  protected postClickOutsideToIframe?: (event: MouseEvent) => void;
  protected changeTheme?: (theme: RemarkTheme) => void;
  protected canScroll = true;

  constructor(props: Props) {
    super(props);
  }

  render() {
    return (
      <div className={`remark42 ${this.props.className || ''}`} id={this.props.id} ref={this.ref} />
    );
  }

  shouldComponentUpdate(nextProps: Props) {
    if (nextProps.theme !== this.props.theme) return true;
    return false;
  }

  componentDidMount() {
    const COMMENT_NODE_CLASSNAME_PREFIX = 'remark42__comment-';

    const remark_config: {
      host: string;
      site_id: string;
      page_title: string;
      url?: string;
      theme: RemarkTheme;
      locale: string;
    } = {
      host: this.props.host || 'https://remark42.radio-t.com',
      site_id: this.props.site_id,
      page_title: this.props.page_title,
      theme: this.props.theme || 'light',
      locale: this.props.locale,
    };

    if (!remark_config.site_id) {
      console.error('Remark42: Site ID is undefined.');
      return;
    }

    remark_config.url = (this.props.url || window.location.href).split('#')[0];

    const node = this.ref.current;

    if (node === null) {
      return null;
    }

    const query = Object.keys(remark_config)
      .map(
        key =>
          `${encodeURIComponent(key)}=${encodeURIComponent(
            (remark_config as Record<string, string>)[key]
          )}`
      )
      .join('&');

    node.innerHTML = `
    <iframe
      src="${remark_config.host}/web/iframe.html?${query}"
      width="100%"
      frameborder="0"
      allowtransparency="true"
      scrolling="no"
      tabindex="0"
      title="Remark42"
      style="width: 1px !important; min-width: 100% !important; border: none !important; overflow: hidden !important;"
      horizontalscrolling="no"
      verticalscrolling="no"
    ></iframe>
  `;

    const iframe = node.getElementsByTagName('iframe')[0];

    this.receiveMessages = function(event) {
      try {
        const data = typeof event.data === 'string' ? JSON.parse(event.data) : event.data;
        if (data.remarkIframeHeight) {
          iframe.style.height = `${data.remarkIframeHeight}px`;
          if (this.canScroll && !data.scrollTo && window.location.hash === '#comments') {
            this.canScroll = false;
            window.scrollTo(
              window.pageXOffset,
              iframe.getBoundingClientRect().top + window.pageYOffset
            );
          }
        }

        if (data.scrollTo) {
          window.scrollTo(
            window.pageXOffset,
            data.scrollTo + iframe.getBoundingClientRect().top + window.pageYOffset
          );
        }

        if (data.hasOwnProperty('isUserInfoShown')) {
          if (data.isUserInfoShown) {
            userInfo.init(data.user || {});
          } else {
            userInfo.close();
          }
        }
      } catch (e) {
        //
      }
    };

    this.postHashToIframe = function(e) {
      const hash = e ? `#${e.newURL.split('#')[1]}` : window.location.hash;

      if (hash.indexOf(`#${COMMENT_NODE_CLASSNAME_PREFIX}`) === 0) {
        if (e) e.preventDefault();

        iframe.contentWindow && iframe.contentWindow.postMessage(JSON.stringify({ hash }), '*');
      }
    };

    this.postClickOutsideToIframe = function(e) {
      if (!iframe.contains(e.target as HTMLElement)) {
        iframe.contentWindow &&
          iframe.contentWindow.postMessage(JSON.stringify({ clickOutside: true }), '*');
      }
    };

    this.changeTheme = function(theme) {
      iframe.contentWindow && iframe.contentWindow.postMessage(JSON.stringify({ theme }), '*');
    };

    window.addEventListener('message', this.receiveMessages);
    window.addEventListener('hashchange', this.postHashToIframe);
    document.addEventListener('click', this.postClickOutsideToIframe);
    setTimeout(this.postHashToIframe, 1000);

    const remarkRootId = 'remark-km423lmfdslkm34';
    const userInfo: {
      node: HTMLElement | null;
      back: HTMLElement | null;
      closeEl: HTMLElement | null;
      iframe: HTMLElement | null;
      style: HTMLElement | null;
      init: (user: Record<string, string>) => void;
      close: () => void;
      delay: number | null;
      events: string[];
      animationStop(): void;
      onAnimationClose: () => void;
      remove: () => void;
      onKeyDown: (e: KeyboardEvent) => void;
    } = {
      node: null,
      back: null,
      closeEl: null,
      iframe: null,
      style: null,
      init(user) {
        this.animationStop();
        if (!this.style) {
          this.style = document.createElement('style');
          this.style.setAttribute('rel', 'stylesheet');
          this.style.setAttribute('type', 'text/css');
          this.style.innerHTML = `
          #${remarkRootId}-node {
            position: fixed;
            top: 0;
            right: 0;
            bottom: 0;
            width: 400px;
            transition: transform 0.4s ease-out;
            max-width: 100%;
            transform: translate(400px, 0);
            z-index: 1036;
          }
          #${remarkRootId}-node[data-animation] {
            transform: translate(0, 0);
          }
          #${remarkRootId}-back {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0,0,0,0.7);
            opacity: 0;
            transition: opacity 0.4s ease-out;
            z-index: 1035;
          }
          #${remarkRootId}-back[data-animation] {
            opacity: 1;
          }
          #${remarkRootId}-close {
            top: 0px;
            right: 400px;
            position: absolute;
            text-align: center;
            font-size: 25px;
            cursor: pointer;
            color: white;
            border-color: transparent;
            border-width: 0;
            padding: 0;
            margin-right: 4px;
            background-color: transparent;
            z-index: 1036;
          }
          @media all and (max-width: 430px) {
            #${remarkRootId}-close {
              right: 0px;
              font-size: 20px;
              color: black;
            }
          }
        `;
        }
        if (!this.node) {
          this.node = document.createElement('div');
          this.node.id = `${remarkRootId}-node`;
        }
        if (!this.back) {
          this.back = document.createElement('div');
          this.back.id = `${remarkRootId}-back`;
          this.back.onclick = () => this.close();
        }
        if (!this.closeEl) {
          this.closeEl = document.createElement('button');
          this.closeEl.id = `${remarkRootId}-close`;
          this.closeEl.innerHTML = '&#10006;';
          this.closeEl.onclick = () => this.close();
        }
        const queryUserInfo =
          `${query}&page=user-info&` +
          `&id=${user.id}&name=${user.name}&picture=${user.picture ||
            ''}&isDefaultPicture=${user.isDefaultPicture || 0}`;
        this.node.innerHTML = `
      <iframe
        src="${remark_config.host}/web/iframe.html?${queryUserInfo}"
        width="100%"
        height="100%"
        frameborder="0"
        allowtransparency="true"
        scrolling="no"
        tabindex="0"
        title="Remark42"
        verticalscrolling="no"
        horizontalscrolling="no"
					/>
				`;
        this.iframe = this.node.querySelector('iframe');
        this.node.appendChild(this.closeEl);
        document.body.appendChild(this.style);
        document.body.appendChild(this.back);
        document.body.appendChild(this.node);
        document.addEventListener('keydown', this.onKeyDown);
        window.setTimeout(() => {
          this.back!.setAttribute('data-animation', '');
          this.node!.setAttribute('data-animation', '');
          this.iframe!.focus();
        }, 400);
      },
      close() {
        if (this.node) {
          this.onAnimationClose();
          this.node.removeAttribute('data-animation');
        }
        if (this.back) {
          this.back.removeAttribute('data-animation');
        }
        document.removeEventListener('keydown', this.onKeyDown);
      },
      delay: null,
      events: ['', 'webkit', 'moz', 'MS', 'o'].map(prefix =>
        prefix ? `${prefix}TransitionEnd` : 'transitionend'
      ),
      onAnimationClose() {
        const el = this.node;
        if (!this.node) {
          return;
        }
        this.delay = window.setTimeout(this.animationStop, 1000);
        this.events.forEach((event: string) =>
          el!.addEventListener(event, this.animationStop, false)
        );
      },
      onKeyDown(e) {
        // ESCAPE key pressed
        if (e.keyCode == 27) {
          userInfo.close();
        }
      },
      animationStop() {
        const t = userInfo;
        if (!t.node) {
          return;
        }
        if (t.delay) {
          clearTimeout(t.delay);
          t.delay = null;
        }
        t.events.forEach(event => t.node!.removeEventListener(event, t.animationStop, false));
        return t.remove();
      },
      remove() {
        const t = userInfo;
        t.node && t.node.remove();
        t.back && t.back.remove();
        t.style && t.style.remove();
      },
    };
  }

  componentWillUnmount() {
    window.removeEventListener('message', this.receiveMessages!);
    window.removeEventListener('hashchange', this.postHashToIframe!);
    document.removeEventListener('click', this.postClickOutsideToIframe!);
  }

  componentDidUpdate(prevProps: Props) {
    if (prevProps.theme !== this.props.theme) {
      this.changeTheme && this.changeTheme(this.props.theme || 'light');
    }
  }
}
