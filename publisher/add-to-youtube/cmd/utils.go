package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func download(from, to string) error {
	log.Infof("Downloading file `%s`to `%s`\n", from, to)

	resp, err := http.Get(from)
	if err != nil {
		return errors.Wrap(err, "Error downloading an audio file")
	}
	defer resp.Body.Close()

	media, err := os.Create(to)
	if err != nil {
		return errors.Wrap(err, "Error creation a file")
	}
	defer media.Close()

	if _, err := io.Copy(media, resp.Body); err != nil {
		return errors.Wrap(err, "Error saving a file to fs")
	}

	log.Infof("File `%s` downloaded\n", to)
	return nil
}

func getEpisodeInfo(id string) (*entry, error) {
	u := fmt.Sprintf("https://radio-t.com/site-api/podcast/%s", id)
	log.Infof("Calling API method `%s`\n", u)
	resp, err := http.Get(u)
	if err != nil {
		return nil, errSiteAPIRequest(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading request body")
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
	log.Infof("Data recieved")
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
	log.Infof("Constructing an episode description")
	return fmt.Sprintf(descriptionFormat, e.URL, id, e.ShowNotes, id, e.AudioURL)
}
