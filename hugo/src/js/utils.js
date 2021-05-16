export function getLocalStorage(key, selector = (s) => s) {
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
  return /^[0,2-9]?[1]$/.test(value)
    ? units[0]
    : /^[0,2-9]?[2-4]$/.test(value)
      ? units[1]
      : units[2];
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
  return new Date(isNaN(time) ? 0 : time * 1000).toISOString().substr(11, 8);
}

export function getTextSnippet(html) {
  const LENGTH = 120;
  const tmp = document.createElement('div');
  tmp.innerHTML = html.replace('</p><p>', ' ').replace(/src=".*"/, '');

  const result = tmp.innerText || '';
  const snippet = result.substr(0, LENGTH);

  return snippet.length === LENGTH && result.length !== LENGTH ? `${snippet}...` : snippet;
}

//https://stackoverflow.com/questions/400212/how-do-i-copy-to-the-clipboard-in-javascript
function fallbackCopyTextToClipboard(text) {
  var textArea = document.createElement("textarea");
  textArea.value = text;

  // Avoid scrolling to bottom
  textArea.style.top = "0";
  textArea.style.left = "0";
  textArea.style.position = "fixed";

  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    var successful = document.execCommand('copy');
    var msg = successful ? 'successful' : 'unsuccessful';
    alert('Временная метка скопирована в буфер обмена', msg, '\nTO DO: сделать красивое оповещение');
  } catch (err) {
    alert('ОШИБКА: невозможно скопировать временную метку в буфер обмена    ', err, "\nTO DO: сделать красивое оповещение");
  }

  document.body.removeChild(textArea);
}

export function copyTextToClipboard(text) {
  if (!navigator.clipboard) {
    fallbackCopyTextToClipboard(text);
    return;
  }
  navigator.clipboard.writeText(text).then(function () {
    alert('Временная метка скопирована в буфер обмена\nTO DO: сделать красивое оповещение');
  }, function (err) {
    alert('ОШИБКА: невозможно скопировать временную метку в буфер обмена    ', err, "\nTO DO: сделать красивое оповещение");
  });
}

export function addTimeToURL(podcastPathname, currentTime) {
  // add full link to URL
  if (window.history.pushState) {
    window.history.pushState(null, null, [podcastPathname, "?t=", currentTime].join(""));
  }
}
