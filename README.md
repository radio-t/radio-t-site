# Радио-Т hugo, скрипты для создания и доставки

## генерация сайта
```
    docker-compose run hugo
```

## скрипты публикации подкаста

- `publisher/make_new_episode.sh` — создает шаблон нового выпуска, берет темы с news.radio-t.com.
– `publisher/make_new_prep.sh` — создает "Темы для ..." следующего выпуска.
- `publisher/upload_mp3.sh` — загружает подкаст во все места, предварительно добавляет mp3 теги и картинку и потом разносит по нодам через внешний ansible контейнер.
- `publisher/deploy.sh` — добавляет в гит и запускает pull + build на мастер. После этого строит лог чата и очищает темы.


## переменные окружения

- `RT_NEWS_ADMIN` user:passwd для news
– `PODCAST_ARCHIVE_CREDS` user:passwd для sftp архивов
