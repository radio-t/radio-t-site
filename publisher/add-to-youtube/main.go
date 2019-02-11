package main

import (
	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/cmd"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		ForceColors:            true,
		DisableLevelTruncation: true,
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	cmd.Execute()
}
