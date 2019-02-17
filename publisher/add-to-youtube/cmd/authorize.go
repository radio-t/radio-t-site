package cmd

import (
	"github.com/pkg/errors"
	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/youtube"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Authorize an user at YouTube",
	Long:  `Authorize an user at YouTube.`,
	Run: func(cmd *cobra.Command, args []string) {
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

	c, err := youtube.New(config, tokenPath, "")
	if err != nil {
		return errors.Wrap(err, "Error creation a youtube client")
	}

	if err := c.Authorize(); err != nil {
		return errors.Wrap(err, "Error authorizing a user in YouTube")
	}
	return nil
}
