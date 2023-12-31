package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/pkg/errors"
	"github.com/umputun/go-flags"

	"github.com/radio-t/publisher/cmd"
)

var opts struct {
	SiteAPI string `long:"site-api" env:"SITE_API" default:"https://radio-t.com/site-api" description:"site API url"`

	NewShowCmd struct {
		NewsAPI string `long:"news" env:"NEWS" default:"https://news.radio-t.com/api/v1/news" description:"news API url"`
		NewsHrs int    `long:"news-hrs" env:"NEWS_HRS" default:"12" description:"news duration in hours"`
		Dest    string `long:"dest" env:"DEST" default:"./content/posts" description:"path to posts"`
	} `command:"new" description:"make new podcast"`

	PrepShowCmd struct {
		Dest string `long:"dest" env:"DEST" default:"./content/posts" description:"path to posts"`
	} `command:"prep" description:"make new prep podcast post"`

	UploadCmd struct {
		Location string `long:"location" env:"LOCATION" default:"/Volumes/Podcasts/radio-t/" description:"podcast location"`
	} `command:"upload" description:"upload podcast"`

	DeployCmd struct {
		NewsAPI    string `long:"news" env:"NEWS" default:"https://news.radio-t.com/api/v1/news" description:"news API url"`
		NewsHrs    int    `long:"news-hrs" env:"NEWS_HRS" default:"12" description:"news duration in hours"`
		NewsPasswd string `long:"news-passwd" env:"NEWS_PASSWD" required:"true" description:"news api admin passwd"`
	} `command:"deploy" description:"deploy podcast to site"`

	Episode int  `short:"e" long:"episode" default:"-1" description:"episode number"`
	Dry     bool `long:"dry" description:"dry run"`
	Dbg     bool `long:"dbg" env:"DEBUG" description:"debug mode"`
}

var revision = "local"

func main() {
	fmt.Printf("rt-publisher, version %s\n", revision)

	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	setupLog(opts.Dbg)

	episodeNum, err := episode()
	if err != nil {
		log.Fatalf("[ERROR] can't get last podcast number, %v", err)
	}
	log.Printf("[DEBUG] episode %d", episodeNum)

	if p.Active != nil && p.Command.Find("new") == p.Active {
		runNew(episodeNum)
	}

	if p.Active != nil && p.Command.Find("prep") == p.Active {
		runPrep(episodeNum)
	}

	if p.Active != nil && p.Command.Find("upload") == p.Active {
		runUpload(episodeNum)
	}

	if p.Active != nil && p.Command.Find("deploy") == p.Active {
		runDeploy(episodeNum)
	}
}

// episode gets the next episode number by hitting site-api
func episode() (int, error) {
	if opts.Episode > 0 {
		return opts.Episode, nil
	}
	client := http.Client{Timeout: 10 * time.Second}
	lastEpisode, err := cmd.LastShow(client, opts.SiteAPI)
	if err != nil {
		return 0, errors.Wrap(err, "can't get last episode")
	}
	return lastEpisode + 1, nil
}

func runNew(episodeNum int) {
	log.Printf("[INFO] make new episode %d", episodeNum)
	prep := cmd.Prep{
		Client:  http.Client{Timeout: 10 * time.Second},
		NewsAPI: opts.NewShowCmd.NewsAPI,
		NewsHrs: opts.NewShowCmd.NewsHrs,
		Dest:    opts.NewShowCmd.Dest,
		Dry:     opts.Dry,
	}
	if err := prep.MakeShow(episodeNum); err != nil {
		log.Fatalf("[ERROR] failed to make new podcast #%d, %v", episodeNum, err)
	}
	log.Printf("[INFO] created new podcast #%d", episodeNum)
	fmt.Printf("%s/podcast-%d.md", opts.PrepShowCmd.Dest, episodeNum) // don't delete! used by external callers

}

func runPrep(episodeNum int) {
	log.Printf("[INFO] prepare next episode post %d", episodeNum)
	prep := cmd.Prep{
		Client: http.Client{Timeout: 10 * time.Second},
		Dest:   opts.PrepShowCmd.Dest,
		Dry:    opts.Dry,
	}
	if err := prep.MakePrep(episodeNum); err != nil {
		log.Fatalf("[ERROR] failed to make new prep #%d, %v", episodeNum, err)
	}
	log.Printf("[INFO] created new prep #%d", episodeNum)
	fmt.Printf("%s/prep-%d.md", opts.PrepShowCmd.Dest, episodeNum) // don't delete! used by external callers
}

func runUpload(episodeNum int) {
	upload := cmd.Upload{
		Executor: &cmd.ShellExecutor{Dry: opts.Dry},
		Location: opts.UploadCmd.Location,
	}
	upload.Do(episodeNum)
	log.Printf("[INFO] deployed #%d", episodeNum)
}

func runDeploy(episodeNum int) {
	deploy := cmd.Deploy{
		Client:     http.Client{Timeout: 10 * time.Second},
		Executor:   &cmd.ShellExecutor{Dry: opts.Dry},
		NewsAPI:    opts.DeployCmd.NewsAPI,
		NewsPasswd: opts.DeployCmd.NewsPasswd,
		NewsHrs:    opts.DeployCmd.NewsHrs,
		Dry:        opts.Dry,
	}
	if err := deploy.Do(episodeNum); err != nil {
		log.Fatalf("[ERROR] failed to deploy #%d, %v", episodeNum, err)
	}
	log.Printf("[INFO] deployed #%d", episodeNum)
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}
