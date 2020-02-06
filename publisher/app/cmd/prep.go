package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

var newShowTmpl = `+++
title = "Радио-Т {{.EpisodeNum}}"
date = {{.TS.Format "2006-01-02T15:04:05"}}
categories = ["podcast"]
image = "https://radio-t.com/images/radio-t/rt{{.EpisodeNum}}.jpg"
filename = "rt_podcast{{.EpisodeNum}}"
+++

![](https://radio-t.com/images/radio-t/rt{{.EpisodeNum}}.jpg)

{{.News}}

*Спонсор этого выпуска [DigitalOcean](https://www.digitalocean.com)*
[аудио](https://cdn.radio-t.com/rt_podcast{{.EpisodeNum}}.mp3) • [лог чата](https://chat.radio-t.com/logs/radio-t-{{.EpisodeNum}}.html)
<audio src="https://cdn.radio-t.com/rt_podcast{{.EpisodeNum}}.mp3" preload="none"></audio>
`

var prepShowTmpl = `+++
title = "Темы для {{.EpisodeNum}}"
date = {{.TS.Format "2006-01-02T15:04:05"}}
categories = ["prep"]
+++
`

// Prep implements both preparation of md file for the new podcast and for prep-show post
type Prep struct {
	Client  http.Client
	NewsAPI string
	NewsHrs int
	Dest    string
	Dry     bool

	now func() time.Time
}

// MakeShow creates md file like podcast-123.md based on newShowTmpl and populated from news response
func (p *Prep) MakeShow(episodeNum int) (err error) {
	if p.now == nil {
		p.now = time.Now
	}

	tp := struct {
		EpisodeNum int
		TS         time.Time
		News       string
	}{
		EpisodeNum: episodeNum,
		TS:         p.now(),
	}

	if tp.News, err = p.lastNews(p.NewsHrs); err != nil {
		return errors.Wrap(err, "failed to load last news")
	}

	return p.applyTemplate(fmt.Sprintf("%s/podcast-%d.md", p.Dest, episodeNum), newShowTmpl, tp)
}

// MakePrep creates a post for news collection, i.e. prep-123.md
func (p *Prep) MakePrep(episodeNum int) (err error) {
	if p.now == nil {
		p.now = time.Now
	}

	tp := struct {
		EpisodeNum int
		TS         time.Time
	}{
		EpisodeNum: episodeNum,
		TS:         p.now(),
	}

	return p.applyTemplate(fmt.Sprintf("%s/prep-%d.md", p.Dest, episodeNum), prepShowTmpl, tp)
}

// applyTemplate writes the applied template to outFile
func (p *Prep) applyTemplate(outFile string, tmpl string, tp interface{}) error {
	t, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return errors.Wrapf(err, "can't parse template")
	}
	msg := bytes.Buffer{}
	if err = t.Execute(&msg, tp); err != nil {
		return errors.Wrapf(err, "can't apply template")
	}
	if p.Dry {
		log.Printf(msg.String())
		return nil
	}
	return errors.Wrapf(ioutil.WriteFile(outFile, msg.Bytes(), 0660), "can't write %s", outFile)
}

// lastNews gets news from news API for the lase hrs hours
func (p *Prep) lastNews(hrs int) (string, error) {
	resp, err := p.Client.Get(fmt.Sprintf("%s/lastmd/%d", p.NewsAPI, hrs))
	if err != nil {
		return "", errors.Wrap(err, "can't get news")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("invalid status code %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "can't read news body")
	}
	return string(b), nil
}
