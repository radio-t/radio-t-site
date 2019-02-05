package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func download(from, to string) error {
	fmt.Printf("Downloading a file `%s` → `%s`\n", from, to)

	resp, err := http.Get(from)
	if err != nil {
		return fmt.Errorf("Error downloading an audio file, got: %v", err)
	}
	defer resp.Body.Close()

	media, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("Error creation a file, got: %v", err)
	}
	defer media.Close()

	if _, err := io.Copy(media, resp.Body); err != nil {
		return fmt.Errorf("Error saving a file to fs, got: %v", err)
	}

	fmt.Printf("File `%s` downloaded\n", to)
	return nil
}

func getEpisodeInfo(id string) (*entry, error) {
	u := fmt.Sprintf("https://radio-t.com/site-api/podcast/%s", id)
	fmt.Printf("Calling API method `%s`\n", u)
	resp, err := http.Get(u)
	if err != nil {
		return nil, errSiteAPIRequest(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading request body, got: %v", err)
	}
	if resp.StatusCode != 200 {
		var e siteAPIError
		if err := json.Unmarshal(b, &e); err != nil {
			return nil, errJSONUnmarshal(err)
		}
		return nil, errSiteAPIRequest(e)
	}
	var e entry
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, errJSONUnmarshal(err)
	}
	return &e, nil
}

const descriptionFormat = `%s

Темы %s выпуска:

%s


Лог чата: http://chat.radio-t.com/logs/radio-t-%s.html
Аудио: %s

Информация о подкасте: https://radio-t.com/info/
Лицензия: https://radio-t.com/license/`

func makeEpisodeDescription(id string, e *entry) string {
	return fmt.Sprintf(descriptionFormat, e.URL, id, e.ShowNotes, id, e.AudioURL)
}
