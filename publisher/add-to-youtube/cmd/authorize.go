package cmd

import (
	"github.com/pkg/errors"
	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/youtube"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	yt "google.golang.org/api/youtube/v3"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Authorize an user at YouTube",
	Long:  `Authorize an user at YouTube.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Running authorize command")
		if err := authorize(); err != nil {
			log.Trace(errors.Wrap(err, "Error authorize command"))
		}
	},
}

func init() {
	rootCmd.AddCommand(authorizeCmd)
}

func authorize() error {
	tokenPath, err := getTokenPath()
	if err != nil {
		return err
	}

	ytConfig := &youtube.Config{TokenPath: tokenPath, OAuth2: config, Scopes: []string{yt.YoutubeUploadScope}}

	log.Debug(tokenPath)

	c, err := youtube.New(ytConfig)
	if err != nil {
		return errors.Wrap(err, "Error creation a youtube client")
	}

	if err := c.Authorize(); err != nil {
		return errors.Wrap(err, "Error authorizing a user in YouTube")
	}
	return nil
}
