package youtube

import (
	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/client"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

func authorize(config *oauth2.Config, pathToSecrets string) error {

	client, err := client.New(youtube.YoutubeUploadScope, &client.Options{PathToSecrets: pathToSecrets, Config: config})
	if err != nil {
		return errYoutubeClientCreate(err)
	}

	if _, err := youtube.New(client); err != nil {
		return errYoutubeClientCreate(err)
	}

	return nil
}
