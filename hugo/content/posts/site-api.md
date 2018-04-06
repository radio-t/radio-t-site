+++
title = "API для сайта Радио-Т"
date = "2017-11-27T14:20:00"
categories = ["news"]
+++

Я давно собирался такое сделать, но все не доходили руки. Идея такая - добавить минимально-достаточный API к нашему сайту для того, чтоб можно было из ваших гениальных приложений с ним (сайтом) работать по людски. И сейчас есть энтузиасты которые парсят страницы, ходят по ссылкам и занимаются прочими непотребными делами, и это не их вина - никакого другого способа просто не было. А теперь такой есть.

<!--more-->

API простой, прямой и явный. Базовый URL `https://radio-t.com/site-api`

- `GET /last/{posts}?categories=podcast,prep` - взять последних {posts} в определенных категориях. Категории опциональны
    - пример: `https://radio-t.com/site-api/last/5?categories=podcast` вернет 5 самых свежих подкастов
    
- `GET /search?q=text-to-search&skip=10&limit=5` - поискать по слову в описании подкаста, `skip` и `limit` опциональны
    - пример: `https://radio-t.com/site-api/search?q=mongo&limit=10` вернет до 10 самых свежих подкастов в описании которых есть слово "mongo"

Оба вызова возвращают JSON лист, с элементами вида:

```go
type Entry struct {
  URL        string      `json:"url"`                   // url поста
  Title      string      `json:"title"`                 // заголовок поста
  Date       time.Time   `json:"date"`                  // дата-время поста в RFC3339 
  Categories []string    `json:"categories"`            // список категорий, массив строк
  Image      string      `json:"image,omitempty"`       // url картинки
  FileName   string      `json:"file_name,omitempty"`   // имя файла
  Body       string      `json:"body,omitempty"`        // тело поста в HTML
  ShowNotes  string      `json:"show_notes,omitempty"`  // пост в текстовом виде
  AudioURL   string      `json:"audio_url,omitempty"`   // url аудио файла
  TimeLabels []TimeLabel `json:"time_labels,omitempty"` // массив временых меток тем
}

type TimeLabel struct {
  Topic    string    `json:"topic"`               // название темы
  Time     time.Time `json:"time"`                // время начала в RFC3339
  Duration int       `json:"duration,omitempty"`  // длительность в секундах
}
```

Баг отчеты, предложения по улучшению/рассширению приветствуются в комментариях.
