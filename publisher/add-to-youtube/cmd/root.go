package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/youtube"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func addToYoutube(id string) error {
	e, err := getEpisodeInfo(id)
	if err != nil {
		return err
	}
	d := makeEpisodeDescription(id, e)

	// prepare temprorary directory
	dir, err := ioutil.TempDir("", "add-to-youtube")
	if err != nil {
		return fmt.Errorf("Error creation a temprorary directory, got: %v", err)
	}
	defer os.RemoveAll(dir)

	filename := path.Join(dir, e.FileName+".mp3")

	if err := download(e.AudioURL, filename); err != nil {
		return err
	}

	ytConfig, err := getConfig()
	if err != nil {
		return err
	}

	c, err := youtube.New(ytConfig)
	if err != nil {
		return err
	}

	v, err := c.Upload(filename, e.Title, d, "22", "podcast,radio-t", "public")
	if err != nil {
		return errJSONUnmarshal(err)
	}

	fmt.Printf("A podcast episode %s uploaded.\nhttps://youtu.be/%s\n", id, v.Id)
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "add-to-youtube {episodeID}",
	Short: "Upload a radio-t podcast episode to Youtube",
	Long: `Upload a radio-t podcast episode to Youtube.

This application is a tool to generate a video file from an audio file,
then uses metadatas from site api to upload it to Youtube.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// if _, err := os.Stat(path.Join(viper.GetString("secrets"), "youtube-secret.json")); os.IsNotExist(err) {
		// 	fmt.Println("Before use this command you need authorize an user")
		// 	os.Exit(1)
		// }

		episodeID := args[0]
		if _, err := strconv.Atoi(episodeID); err != nil {
			fmt.Printf("{episodeID} must be a number, got: %v\n", err)
			os.Exit(1)
		}

		if err := addToYoutube(episodeID); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.add-to-youtube.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".add-to-youtube" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".add-to-youtube")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
