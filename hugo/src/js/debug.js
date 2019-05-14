let debug;

if (process.env.NODE_ENV === 'production') {
  debug = () => () => {};
} else {
  debug = require('debug');
}

const pool = {};

export default function (component, ...args) {
  pool[component] = pool[component] || debug(`hugo:${component}`);
  if (args.length) return pool[component](...args);
  return pool[component];
}
