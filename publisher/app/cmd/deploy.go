package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

//go:generate moq --out mocks/executor.go --pkg mocks --with-resets --skip-ensure . Executor

// Deploy delivers site update
type Deploy struct {
	Executor
	NewsPasswd string
	NewsAPI    string
	NewsHrs    int
	Client     http.Client
	Dry        bool
}

var superUsersTelegram = []string{"umputun", "bobuk", "ksenks", "grayodesa", "aleks_sys"}

// Do performs a series of actions to deploy a new episode.
// It takes an episode number as input and returns an error if any of the actions fail.
// It performs the following actions:
//  1. Commit the new episode to git.
//  2. Update the remote hugo site via ssh.
//  3. Create the chat log.
//  4. Archive the news.
func (d *Deploy) Do(episodeNum int) error {
	log.Printf("[INFO] commit new episode to git")
	d.Run(fmt.Sprintf(`git pull && git commit -am "episode %d" && git push`, episodeNum))

	log.Printf("[INFO] remote site update")
	d.Run("ssh umputun@master.radio-t.com",
		`"cd /srv/site.hugo && git pull && /usr/local/bin/docker-compose run --rm hugo"`)

	log.Printf("[INFO] create chat log")
	slParams := []string{}
	for _, s := range superUsersTelegram {
		slParams = append(slParams, fmt.Sprintf("--super=%s", s))
	}
	d.Run("ssh umputun@master.radio-t.com", fmt.Sprintf(`"docker exec -i super-bot /srv/telegram-rt-bot %s --dbg --export-num=%d --export-path=/srv/html"`,
		strings.Join(slParams, " "), episodeNum))

	log.Printf("[INFO] archive news")
	err := d.archiveNews()
	return err
}

// archiveNews invokes news-api like https://news.radio-t.com/api/v1/news/active/last/12
func (d *Deploy) archiveNews() error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/active/last/%d", d.NewsAPI, d.NewsHrs), http.NoBody)
	if err != nil {
		return errors.Wrap(err, "failed to prepare news archive request")
	}
	if d.Dry {
		log.Printf("[INFO] %s", req.URL.String())
		return nil
	}

	req.SetBasicAuth("admin", d.NewsPasswd)
	resp, err := d.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "can't make news archive request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "news archive request returned %s", resp.Status)
	}
	return nil
}
