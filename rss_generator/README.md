# построитель RSS-фида для подкастов

Строит RSS-фиды для основного и архвиного подкастов. Эта программа собирается в hugo контейнере и запускается внутри него при сборке сайта.

## Как работает

- загружает все посты из hugo/content/posts, фильтрует по тегу `podcast`, сортирует по дате
- парсит метаданные из каждого поста
- строит RSS-фиды для основного и архивного подкастов применяя template из `rss_template.gp` и метаданные из постов
- для определения длины подкаста использует берет content-length из ответа сервера на HEAD запрос к mp3 файлу подкаста
- для определения даты публикации подкаста использует дату из метаданных поста


```
Usage:
  rss_generator [OPTIONS]

Application Options:
      --hugo-posts= directory of hugo posts (default: ./content/posts)
      --save-to=    directory for generated feeds (default: /srv/hugo/public)
      --dry         dry run
      --dbg         debug mode [$DEBUG]

Help Options:
  -h, --help        Show this help message
```

Вызывается автоматически при сборке сайта hugo контейнером, часть `exec.sh`. Все параметры передаются через переменные окружения внутри контейнера.