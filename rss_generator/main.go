package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-pkgz/lgr"
	"github.com/umputun/go-flags"

	"github.com/radio-t/radio-t-site/rss_generator/proc"
)

var opts struct {
	HugoPosts string `long:"hugo-posts" default:"./content/posts" description:"directory of hugo posts"`
	SaveTo    string `long:"save-to" default:"/srv/hugo/public" description:"directory for generated feeds"`
	Dry       bool   `long:"dry" description:"dry run"`
	Dbg       bool   `long:"dbg" env:"DEBUG" description:"debug mode"`
}

var feeds = []proc.FeedConfig{
	{
		Name:            "podcast",
		Title:           "Радио-Т",
		Image:           "https://radio-t.com/images/covers/cover.png",
		Count:           20,
		Size:            true,
		FeedSubtitle:    "Подкаст выходного дня - импровизации на темы высоких технологий",
		FeedDescription: "Разговоры на темы хайтек, высоких компьютерных технологий, гаджетов, облаков, программирования и прочего интересного из мира ИТ.",
		Verbose:         true,
	},
	{
		Name:            "podcast-archives",
		Title:           "Радио-Т Архивы",
		Image:           "https://radio-t.com/images/covers/cover-archive.png",
		Count:           1000,
		Size:            false,
		FeedSubtitle:    "Подкаст выходного дня - импровизации на темы высоких технологий",
		FeedDescription: "Разговоры на темы хайтек, высоких компьютерных технологий, гаджетов, облаков, программирования и прочего интересного из мира ИТ.",
	},
	{
		Name:            "podcast-archives-short",
		Title:           "Радио-Т Архивы",
		Image:           "https://radio-t.com/images/covers/cover-archive.png",
		Count:           200,
		Size:            false,
		FeedSubtitle:    "Подкаст выходного дня - импровизации на темы высоких технологий",
		FeedDescription: "Разговоры на темы хайтек, высоких компьютерных технологий, гаджетов, облаков, программирования и прочего интересного из мира ИТ.",
	},
}

var revision = "local"

func main() {
	fmt.Printf("rt-feed-generator, version %s\n", revision)

	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}
	setupLog(opts.Dbg)

	hs := proc.HugoSite{Location: opts.HugoPosts, BaseURL: "https://radio-t.com"}
	posts, err := hs.ReadPosts()
	if err != nil {
		log.Fatalf("error reading posts from %s: %v", opts.HugoPosts, err)
	}

	g := proc.RSSGenerator{
		Client:          &http.Client{Timeout: 10 * time.Second},
		BaseURL:         "https://radio-t.com",
		BaseArchiveURL:  "http://archive.rucast.net/radio-t/media",
		BaseCdnURL:      "http://cdn.radio-t.com",
		RssRootLocation: opts.SaveTo,
	}

	for _, feed := range feeds {
		log.Printf("[INFO] generating feed %q", feed.Name)
		feedData, err := g.MakeFeed(feed, posts)
		if err != nil {
			log.Printf("[WARN] error generating feed data for %s: %v\n", feed.Name, err)
			continue
		}

		if err := g.Save(feed, feedData); err != nil {
			log.Printf("[WARN] error saving feed for %s: %v\n", feed.Name, err)
		}
	}
	log.Printf("[INFO] done generating %d feeds", len(feeds))
}

func setupLog(dbg bool) {
	logOpts := []lgr.Option{lgr.Msec, lgr.LevelBraces, lgr.StackTraceOnError}
	if dbg {
		logOpts = []lgr.Option{lgr.Debug, lgr.CallerFile, lgr.CallerFunc, lgr.Msec, lgr.LevelBraces, lgr.StackTraceOnError}
	}
	lgr.SetupStdLogger(logOpts...)
	lgr.Setup(logOpts...)
}
