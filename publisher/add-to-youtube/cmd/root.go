package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/youtube"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yt "google.golang.org/api/youtube/v3"
)

var (
	cfgFile     string
	pathToToken string
	config      []byte
)

func addToYoutube(id string) error {
	e, err := getEpisodeInfo(id)
	if err != nil {
		return err
	}
	d := makeEpisodeDescription(id, e)

	log.Info("Creating temporary directory")
	dir, err := ioutil.TempDir("", "add-to-youtube")
	if err != nil {
		return errors.Wrap(err,"Error creation a temprorary directory")
	}
	defer func() {
		log.Infof("Removing temporary directory `%s`", dir)
		os.RemoveAll(dir)
	}()

	filename := path.Join(dir, e.FileName+".mp3")

	if err := download(e.AudioURL, filename); err != nil {
		return err
	}

	ytConfig := &youtube.Config{TokenPath: pathToToken, OAuth2: config, Scopes: []string{yt.YoutubeUploadScope}}

	c, err := youtube.New(ytConfig)
	if err != nil {
		return err
	}

	v, err := c.Upload(filename, e.Title, d, "22", "podcast,radio-t", "private")
	if err != nil {
		return errJSONUnmarshal(err)
	}

	log.Infof("A podcast episode %s uploaded.\nhttps://youtu.be/%s\n", id, v.Id)
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "add-to-youtube {episodeID}",
	Short: "Upload a radio-t podcast episode to Youtube",
	Long: `Upload a radio-t podcast episode to Youtube.

This application is a tool to generate a video file from an audio file via ffmpeg,
then uses metadatas from site api to upload it to Youtube.`,
	Args: cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if config, err = getClientSecretJSON(); err != nil {
			log.Fatal(err)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		pathToToken, err = getTokenPath()
		if err != nil {
			log.Fatal(err)
		}
		if _, err := os.Stat(pathToToken); os.IsNotExist(err) {
			log.Fatal("Required user authorization")
		}
		cd, err := os.Getwd()
		if err != nil {
			log.Fatal(errors.WithStack(err))
		}
		coverPath := path.Join(cd, "assets/cover.webp")
		if _, err := os.Stat(coverPath); os.IsNotExist(err) {
			log.Fatalf("An image cover not found at `%s`", coverPath)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Running upload command")
		episodeID := args[0]
		if _, err := strconv.Atoi(episodeID); err != nil {
			log.Fatalf("{episodeID} must be a number, got: %s", err)
		}

		if err := addToYoutube(episodeID); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.add-to-youtube.yaml)")
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
			log.Fatal(err)
		}

		// Search config in home directory with name ".add-to-youtube" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".add-to-youtube")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file: ", viper.ConfigFileUsed())
	}
}
