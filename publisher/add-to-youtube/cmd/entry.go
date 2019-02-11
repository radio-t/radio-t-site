package cmd

import "time"

type entry struct {
	URL        string      `json:"url"`                   // url поста
	Title      string      `json:"title"`                 // заголовок поста
	Date       time.Time   `json:"date"`                  // дата-время поста в RFC3339
	Categories []string    `json:"categories"`            // список категорий, массив строк
	Image      string      `json:"image,omitempty"`       // url картинки
	FileName   string      `json:"file_name,omitempty"`   // имя файла
	Body       string      `json:"body,omitempty"`        // тело поста в HTML
	ShowNotes  string      `json:"show_notes,omitempty"`  // пост в текстовом виде
	AudioURL   string      `json:"audio_url,omitempty"`   // url аудио файла
	TimeLabels []timeLabel `json:"time_labels,omitempty"` // массив временых меток тем
}

type timeLabel struct {
	Topic    string    `json:"topic"`              // название темы
	Time     time.Time `json:"time"`               // время начала в RFC3339
	Duration int       `json:"duration,omitempty"` // длительность в секундах
}

type siteAPIError struct {
	Message string `json:"error"`
}

func (e siteAPIError) Error() string {
	return e.Message
}
