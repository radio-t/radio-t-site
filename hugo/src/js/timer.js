import { getUnits } from './utils';

function padTime(n) {
  return `0${n}`.slice(-2);
}

function formatSeconds(totalSeconds) {
  const seconds = totalSeconds % 60;
  const minutes = Math.round((totalSeconds - seconds) / 60) % 60;
  let hours = Math.round((totalSeconds - seconds - minutes * 60) / 3600);
  const days = (hours - (hours % 24)) / 24;

  hours %= 24;

  let html = '';
  if (days > 0) {
    html += `${days} ${getUnits(days, ['день', 'дня', 'дней'])} `;
  }

  html += `${padTime(hours)}:${padTime(minutes)}<span style="opacity: .5;">:${padTime(
    seconds
  )}</span>`;
  return html;
}

/**
 * Show airs from 20:00 to 23:00 Saturday UTC
 * So we'll skip all the timezone's troubles
 */
export function timer(now = new Date()) {
  const showStart = new Date(
    Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate() + 6 - now.getUTCDay(), 20)
  );

  let totalSeconds = Math.floor((showStart - now) / 1000);
  if (totalSeconds <= -3 * 60 * 60) {
    totalSeconds += 7 * 24 * 60 * 60;
  }

  const isOnline = totalSeconds <= 0;

  const html = isOnline ? 'Мы в эфире!' : formatSeconds(totalSeconds);

  return { isOnline, html };
}
