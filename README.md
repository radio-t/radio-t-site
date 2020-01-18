# Радио-Т hugo, скрипты для создания и доставки

## генерация сайта
```bash
docker-compose build hugo
docker-compose run --rm hugo
```

## публикация подкаста

Перед использованием, необходимо иметь собранный docker `publisher`. Команда сборки при помощи docker-compose (конфиг в руте репозитария): `docker-compose build publisher`.

Скрипты публикации могут быть вызваны при помощи make в директории `./publisher`:

- `make` - список доступных команд
- `make new-episode` — создает шаблон нового выпуска, берет темы с news.radio-t.com
- `make new-prep` — создает шаблон "Темы для ..." следующего выпуска
- `make upload-mp3 FILE=rt_podcast685.mp3` - добавляет mp3 теги и картинку в файл подкаста, после чего разносит по нодам его через внешний ansible контейнер. Для выполнения необходимо подключить в docker-compose конфиге директорию с mp3 файлами подкаста как volume в сервис publisher
- `make deploy` — добавляет в гит и запускает pull + build на мастер. После этого строит лог чата и очищает темы


## переменные окружения

- `RT_NEWS_ADMIN` user:passwd для news
- `PODCAST_ARCHIVE_CREDS` user:passwd для sftp архивов

## фронтенд

```bash
# node 10
cd hugo

npm install

# разработка на localhost:3000
# с hugo LiveReload, без turbolinks
npm run start
# без hugo LiveReload, с turbolinks
npm run start-turbo

# сборка для прода
# результаты сборки:
# - hugo/static/build/
# - hugo/data/manifest.json
npm run production
```

лого в `src/images/`

фавиконки в `static/` и описаны в `layouts/partials/favicons.html`

обложки в `static/images/covers/` (для сохранения совместимости также оставлены обложки `static/images/cover.jpg` и `static/images/cover_rt_big_archive.png`)
