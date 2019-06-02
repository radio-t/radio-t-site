import padStart from 'lodash/padStart';

export function getLocalStorage(key, selector = s => s) {
  let data;
  try {
    data = JSON.parse(localStorage.getItem(key) || '{}');
  } catch (e) {
    //
  }
  return selector(data);
}

export function updateLocalStorage(key, fn) {
  try {
    const newValue = fn(getLocalStorage(key));
    if (typeof newValue === 'undefined') return;
    localStorage.setItem(key, JSON.stringify(newValue));
  } catch (e) {
    //
  }
}

export function getUnits(value, units) {
  return (/^[0,2-9]?[1]$/.test(value)) ? units[0] : ((/^[0,2-9]?[2-4]$/.test(value)) ? units[1] : units[2]);
}

// 00:02:24 => 144
export function parseTime(time) {
  return time
    .split(':')
    .reverse()
    .reduce((acc, curr, i) => acc + parseInt(curr) * Math.pow(60, i), 0);
}

// 144 => 00:02:24
export function composeTime(time) {
  const pieces = [];
  time = parseInt(time);
  while (time) {
    pieces.push(time % 60);
    time = Math.floor(time / 60);
  }
  while (pieces.length < 3) pieces.push(0);
  return pieces.reverse().map((t) => padStart(t, 2, '0')).join(':');
}

export function getTextSnippet(html) {
  const LENGTH = 120;
  const tmp = document.createElement('div');
  tmp.innerHTML = html.replace('</p><p>', ' ');

  const result = tmp.innerText || '';
  const snippet = result.substr(0, LENGTH);

  return snippet.length === LENGTH && result.length !== LENGTH ? `${snippet}...` : snippet;
}
