import hljs from 'highlight.js/lib/highlight';
import 'highlight.js/styles/atom-one-dark.css';

const languages = {
  bash: require('highlight.js/lib/languages/bash'),
  go: require('highlight.js/lib/languages/go'),
  python: require('highlight.js/lib/languages/python'),
  java: require('highlight.js/lib/languages/java'),
  javascript: require('highlight.js/lib/languages/javascript'),
  json: require('highlight.js/lib/languages/json'),
};

Object.keys(languages).forEach((lang) => hljs.registerLanguage(lang, languages[lang]));

export default hljs;
