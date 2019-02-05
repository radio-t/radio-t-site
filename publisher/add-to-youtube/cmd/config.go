package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/youtube"
	yt "google.golang.org/api/youtube/v3"
)

const (
	secrets                       = "ADD_RADIOT_TO_YOUTUBE_SECRETS_PATH"
	addRadioT2YoutubeClientSecret = "ADD_RADIOT_TO_YOUTUBE_CLIENT_SECRET"
)

func getOAuth2Config() (*oauth2.Config, error) {
	s := viper.GetString(addRadioT2YoutubeClientSecret)
	config, err := google.ConfigFromJSON([]byte(s), yt.YoutubeUploadScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret string from ENV variable `%s` or app config to oauth2 config: %v",
			addRadioT2YoutubeClientSecret, err)
	}
	return config, nil
}

func getConfig() (*youtube.Config, error) {
	c, err := getOAuth2Config()
	if err != nil {
		return nil, fmt.Errorf("Unable to get oauth2 config, got : %v", err)
	}

	secretsPath := viper.GetString(secrets)
	if _, err := os.Stat(secretsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Path to directory with secrets from ENV variable %s or app config not exists: %v", secrets, err)
	}

	return &youtube.Config{OAuth2: c, SecretsPath: secretsPath}, nil
}
