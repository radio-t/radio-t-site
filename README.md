# Радио-Т hugo, скрипты для создания и доставки

## генерация сайта

```bash
docker-compose build hugo
docker-compose run --rm hugo
```

## публикация подкаста

Перед использованием, необходимо иметь собранный docker образ `publisher`. Команда сборки при помощи docker-compose (конфиг в руте репозитория): `docker-compose build publisher`.

Скрипты публикации могут быть вызваны при помощи make в директории `./publisher`:

- `make new` — создает шаблон нового выпуска, темы берутся с news.radio-t.com, номер выпуска берется из api сайта
- `make prep` — создает шаблон "Темы для ..." следующего выпуска, номер выпуска берется из api сайта
- `make upload-mp3 EPISODE=685` - добавляет mp3 теги и картинку в файл эпизода подкаста, после чего разносит его по нодам через внешний ansible контейнер. Для выполнения необходимо подключить в docker-compose конфиге директорию с mp3 файлами подкаста как volume в сервис publisher
- `make deploy` — добавляет в гит и запускает push + build на мастер. После этого строит лог чата и очищает темы

## переменные окружения

- `RT_NEWS_ADMIN` user:passwd для news
- `PODCAST_ARCHIVE_CREDS` user:passwd для sftp архивов

## фронтенд

### зависимости

- [Node v22](https://nodejs.org/en/download/package-manager)
- [GoLang](https://go.dev/doc/install)
- [Hugo v0.81.0](https://gohugo.io/installation/macos/#build-from-source)

### девелопмент

```bash
# node 10
cd hugo

npm install

# разработка на localhost:3000
# с hugo LiveReload, без turbolinks
npm run dev
# без hugo LiveReload, с turbolinks
npm run dev:turbo

# сборка для prod
# результаты сборки:
# - hugo/static/build/
# - hugo/data/manifest.json
npm run production
```

### файловая структура

- лого в `src/images/`
- favicons в `static/` и описаны в `layouts/partials/favicons.html`
- обложки в `static/images/covers/` (для сохранения совместимости также оставлены обложки `static/images/cover.jpg` и `static/images/cover_rt_big_archive.png`)
