import startOfWeek from 'date-fns/startOfWeek';
import addWeeks from 'date-fns/addWeeks';
import addMinutes from 'date-fns/addMinutes';
import { getUnits } from './utils';

const showTime = {
  day: 6, // 6=Sat, 0..6 - Sun..Sat
  hours: 23,
  minutes: 0,
};

const durationMinutes = 3 * 60;

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

function inMoscow(date) {
  const timeInMoscow = new Date(date.getTime());
  timeInMoscow.setMinutes(timeInMoscow.getMinutes() + timeInMoscow.getTimezoneOffset() + 3 * 60);

  return timeInMoscow;
}

function timeToMinutes(time) {
  return time.day * 24 * 60 + time.hours * 60 + time.minutes;
}

function showEndTime(time, duration) {
  let endTime = startOfWeek(new Date() /*, {weekStartsOn: 0}*/); // week starts on Sunday

  endTime = addMinutes(endTime, timeToMinutes(time) + duration);

  return {
    day: endTime.getDay(),
    hours: endTime.getHours(),
    minutes: endTime.getMinutes(),
  };
}

function getNextOccurrence(now, time) {
  const offsetMinutes = timeToMinutes(time);
  return addMinutes(addWeeks(startOfWeek(addMinutes(now, -offsetMinutes)), 1), offsetMinutes);
}

export function timer(now = new Date()) {
  const timeInMoscow = inMoscow(now);

  const showEnd = getNextOccurrence(timeInMoscow, showEndTime(showTime, durationMinutes));
  const showStart = addMinutes(showEnd, -durationMinutes);

  const totalSeconds = Math.floor((showStart - timeInMoscow) / 1000);
  const isOnline = totalSeconds <= 0;

  const html = isOnline ? 'Мы в эфире!' : formatSeconds(totalSeconds);

  return { isOnline, html };
}
