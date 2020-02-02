# Скрипты публикации страницы нового эпизода на сайт, заливки файла подкаста во все места

## Как пользоваться скриптами в этой директории?

Перед использованием, необходимо собрать docker образ при помощи docker-compose (конфиг в руте репозитария), команда: `docker-compose build publisher`.

После сборки образа, скриптами публикации можно пользоваться как с помощью `make`:

- `make` - список доступных команд
- `make new-episode` — создает шаблон нового выпуска, темы берутся с news.radio-t.com
- `make new-prep` — создает шаблон "Темы для ..." следующего выпуска
- `make upload-mp3 FILE=rt_podcast685/rt_podcast685.mp3` - добавляет mp3 теги и картинку в файл подкаста, после чего разносит его по нодам через внешний ansible контейнер. Для выполнения необходимо подключить в docker-compose конфиге директорию с mp3 файлами подкаста как volume в сервис publisher
- `make deploy` — добавляет в гит и запускает pull + build на мастер. После этого строит лог чата и очищает темы

так и при помощи `docker-compose`:

- `docker-compose run --rm publisher --list` - вывод списка возможных команд для образа
- `docker-compose run --rm publisher --help set-mp3-chapters` - вывод справки по конкретной команде
- `docker-compose run --rm publisher new-episode` — создает шаблон нового выпуска, темы берутся с news.radio-t.com
- `docker-compose run --rm publisher new-prep` — создает шаблон "Темы для ..." следующего выпуска
- `docker-compose run --rm publisher upload-mp3 rt_podcast685/rt_podcast685.mp3` — загружает подкаст во все места, предварительно добавляет mp3 теги и картинку и потом разносит по нодам через внешний ansible контейнер. Для выполнения необходимо подключить в docker-compose конфиге директорию с mp3 файлами подкаста как volume в сервис publisher
- `docker-compose run --rm publisher deploy` — добавляет в гит и запускает push + build на мастер. После этого строит лог чата и очищает темы
 

## Для разработчика

После изменений python скриптов, желательно прогнать следующие линтеры и форматтеры:
 - `isort --lines 119 -y`
 - `black --line-length 120 --target-version py38 ./*/*.py`
 - `flake8 . --max-line-length=120`
