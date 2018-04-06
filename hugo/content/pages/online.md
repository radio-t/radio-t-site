{
   "title": "Слушать Online вещание",
   "url": "/online"
}

![](/images/listen.jpg)

Запись подкаста производится по субботам, в 23:00мск. В это время, вы можете слушать нас в прямом эфире по ссылке: [http://stream.radio-t.com](http://stream.radio-t.com) или прямо тут: <span id="play-stream" class="disabled"></span>
<audio id="stream" src="http://stream.radio-t.com"></audio>

Для обратной связи: [gitter чат](https://gitter.im/radio-t/chat) или кликнуть на "open chat" в нижнем правом углу.

Обратный отсчет: <span id="timer"></span>

<script>
function setShowTimer() {
    function getUnits(value, units) {
        return (/^[0,2-9]?[1]$/.test(value)) ? units[0] : ((/^[0,2-9]?[2-4]$/.test(value)) ? units[1] : units[2])
    }

    var timeInMoscow = new Date();
    timeInMoscow.setMinutes(timeInMoscow.getMinutes() + timeInMoscow.getTimezoneOffset() + 3 * 60);

    var nextShow = new Date(timeInMoscow);
    nextShow.setDate(nextShow.getDate() + 6 - nextShow.getDay());
    nextShow.setHours(23, 0, 0, 0);

    var totalSeconds = Math.floor((nextShow - timeInMoscow) / 1000);

    if (totalSeconds < 0) {
        return "Вещаем!";
    }

    var seconds = totalSeconds % 60,
        minutes = Math.round((totalSeconds - seconds) / 60) % 60,
        hours = Math.round((totalSeconds - seconds - minutes * 60) / 3600),
        days = (hours - hours % 24) / 24;

    hours %= 24;

    var result = "",
        daysList = ['день', 'дня', 'дней'],
        hoursList = ['час', 'часа', 'часов'],
        minutesList = ['минута', 'минуты', 'минут'],
        secondsList = ['секунда', 'секунды', 'секунд'];

    if (days > 0) {
        result += days + ' ' + getUnits(days, daysList) + ' ';
    }

    result += (('0' + hours).slice(-2) + ' ' + getUnits(hours, hoursList) + ' ') +
              (('0' + minutes).slice(-2) + ' ' + getUnits(minutes, minutesList) + ' ') +
              (('0' + seconds).slice(-2) + ' ' + getUnits(seconds, secondsList));

    return result;
}

var t = document.getElementById('timer');

t.textContent = setShowTimer();
window.setInterval(function() {
    t.textContent = setShowTimer();
}, 999);
</script>


<script>
  ((window.gitter = {}).chat = {}).options = {
    room: 'radio-t/chat'
  };
</script>
<script src="https://sidecar.gitter.im/dist/sidecar.v1.js" async defer></script>

<script type="text/javascript">
var playButton = document.getElementById('play-stream'),
  audio = document.getElementById('stream'),
  src = audio.src,
  timer = document.getElementById('timer');
if (playButton) {
var errorHandler = function() {
playButton.classList.add('disabled');
if (timer.textContent == 'Вещаем!') {
setTimeout(function() {
audio.pause();
        audio.src = null;
        audio.src = src;
        audio.play();

        playButton.classList.remove('disabled');
      }, 5000);
    }
  };

  playButton.addEventListener('click', function(e) {
    var target = e.target;

    if (audio.paused) {
      audio.src = src;
      audio.play();
      target.classList.remove('disabled');

      audio.addEventListener('error', errorHandler);
    } else {
      audio.removeEventListener('error', errorHandler);

      audio.pause();
      audio.src = null;
      target.classList.add('disabled');
    }
  });
}
</script>
