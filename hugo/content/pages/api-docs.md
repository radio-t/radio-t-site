{
   "title": "API Радио-Т",
   "url": "/api-docs"
}

#### API сайта radio-t.com

Базовый URL `https://radio-t.com/site-api`

- `GET /last/{posts}?categories=podcast,prep` — взять последних `{posts}` в определенных категориях. Категории опциональны

    пример: `https://radio-t.com/site-api/last/5?categories=podcast` вернет 5 самых свежих подкастов
    
- `GET /search?q=text-to-search&skip=10&limit=5` — поискать по слову в описании подкаста, `skip` и `limit` опциональны
    
    пример: `https://radio-t.com/site-api/search?q=mongo&limit=10` вернет до 10 самых свежих подкастов в описании которых есть слово "mongo"

oба вызова возвращают JSON лист, с элементами вида:

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

- `GET /podcast/{num}` — вернуть информацию о подкасте с заданным номером, возвращает JSON `Entry` 

    пример: `https://radio-t.com/site-api/podcast/223` 


#### API новостей news.radio-t.com

Базовый URL `https://news.radio-t.com/api/v1/`

- `GET /news/active/last/{hrs}` — взять темы активированные в последние `{hrs}` часов. Возвращает массив `Article`
- `GET /news/last/{count}` — возвращает последние добавленные темы
- `GET /news/slug/{slug}` — тема по slug
- `GET /news/domain/#domain` — темы для домена
- `GET /news/active` — возвращет активную, в этот момент, тему
- `GET /news/rss/{count}` — RSS с последнними `{count}` темами
- `GET /show/start` — время начала подкаста


```go
type Article struct {
	Title      string        `json:"title"`     // заголовок темы/новости
	Content    string        `json:"content"`   // полный текст новости
	Snippet    string        `json:"snippet"`   // короткое текстовое описание
	MainPic    string        `json:"pic"`       // ссылка на основную картинку
	Link       string        `json:"link"`      // ссылка на оригинал 
	Author     string        `json:"author"`    // автор новости
	Ts         time.Time     `json:"ts"`        // дата-время оригинла
	AddedTS    time.Time     `json:"ats"`       // дата-время добавления на сайт 
	Active     bool          `json:"active"`    // флаг текущей активности
	ActiveTS   time.Time     `json:"activets"`  // дата-время активации
	Geek       bool          `json:"geek"`      // флаг гиковской темы
	Votes      int           `json:"votes"`     // колличество голосов за тему
	Deleted    bool          `json:"del"`       // флаг удаления 
	Archived   bool          `json:"archived"`  // флаг архивации
	Slug       string        `json:"slug"`      // slug новости
	SourceFeed string        `json:"feed"`      // RSS фид источника 
	Domain     string        `json:"domain"`    // домен новости
	Comments   int           `json:"comments"`  // число комментариев
	Likes      int           `json:"likes"`     // число лайков
}
```
