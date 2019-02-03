package cmd

import (
	"fmt"
	"os"

	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/youtube"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Authorize a user in YouTube",
	Long:  `Authorize a user in YouTube.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := youtube.Authorize(viper.GetString("secrets")); err != nil {
			fmt.Printf("Error authorizing a user in YouTube, got: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(authorizeCmd)
}
