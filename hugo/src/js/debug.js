import debug from 'debug';

const pool = {};

export default function (component, ...args) {
  pool[component] = pool[component] || debug(`hugo:${component}`);
  if (args.length) return pool[component](...args);
  return pool[component];
}
