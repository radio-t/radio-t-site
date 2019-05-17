# Радио-Т hugo, скрипты для создания и доставки

## генерация сайта
```bash
docker-compose build hugo
docker-compose run --rm hugo
```

## скрипты публикации подкаста

- `publisher/make_new_episode.sh` — создает шаблон нового выпуска, берет темы с news.radio-t.com.
- `publisher/make_new_prep.sh` — создает "Темы для ..." следующего выпуска.
- `publisher/upload_mp3.sh` — загружает подкаст во все места, предварительно добавляет mp3 теги и картинку и потом разносит по нодам через внешний ansible контейнер.
- `publisher/deploy.sh` — добавляет в гит и запускает pull + build на мастер. После этого строит лог чата и очищает темы.


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
