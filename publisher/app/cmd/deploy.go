package cmd

import (
	"fmt"
	"net/http"

	log "github.com/go-pkgz/lgr"
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

// Do run deploy sequence for the given episodeNum
// may panic on executor error
func (d *Deploy) Do(episodeNum int) error {
	log.Printf("[INFO] commit new episode to git")
	d.Run(fmt.Sprintf(`git pull && git commit -am "episode %d" && git push`, episodeNum))

	log.Printf("[INFO] remote site update")
	d.Run("ssh umputun@master.radio-t.com", `cd /srv/site.hugo && git pull && docker-compose run --rm hugo`)

	log.Printf("[INFO] create chat log")
	d.Run("ssh umputun@master.radio-t.com", fmt.Sprintf(`docker exec -i super-bot /srv/telegram-rt-bot --super=umputun --super=bobuk --super=ksenks --super=grayodesa --dbg --export-num=%d --export-path=/srv/html`, episodeNum))

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
