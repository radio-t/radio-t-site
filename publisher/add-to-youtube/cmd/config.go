package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	tokenPathKey        = "ADD_RADIOT_TO_YOUTUBE_SECRET_TOKEN_PATH"
	clientSecretJSONKey = "ADD_RADIOT_TO_YOUTUBE_CLIENT_SECRET_JSON"
)

func getClientSecretJSON() ([]byte, error) {
	s := viper.GetString(clientSecretJSONKey)
	if s == "" {
		return nil, errors.Errorf("Missing client secret json string in ENV variable or app config by key=`%s`",
			clientSecretJSONKey)
	}
	return []byte(s), nil
}

func getTokenPath() (string, error) {
	tokenPath := viper.GetString(tokenPathKey)
	if tokenPath == "" {
		return "", errors.Errorf("Missing tokenPath in ENV variable or app config by key=`%s`", tokenPathKey)
	}
	return tokenPath, nil
}
