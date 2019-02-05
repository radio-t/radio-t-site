package cmd

import (
	"fmt"
	"os"

	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/youtube"

	"github.com/spf13/cobra"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Authorize a user in YouTube",
	Long:  `Authorize a user in YouTube.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := authorize(); err != nil {
			fmt.Printf("Error authorize command, got: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(authorizeCmd)
}

func authorize() error {
	ytConfig, err := getConfig()
	if err != nil {
		return fmt.Errorf("Error getting a config, got: %v", err)
	}

	c, err := youtube.New(ytConfig)
	if err != nil {
		return fmt.Errorf("Error creation a youtube client, got: %v", err)
	}

	if err := c.Authorize(); err != nil {
		return fmt.Errorf("Error authorizing a user in YouTube, got: %v", err)
	}
	return nil
}
