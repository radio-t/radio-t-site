package youtube

import (
	"github.com/radio-t/radio-t-site/publisher/podcast-to-youtube/client"
	"google.golang.org/api/youtube/v3"
)

// Authorize authorizes an user in Youtube service.
func Authorize(pathToSecrets string) error {

	client, err := client.New(youtube.YoutubeUploadScope, &client.Options{PathToSecrets: pathToSecrets})
	if err != nil {
		return errYoutubeClientCreate(err)
	}

	if _, err := youtube.New(client); err != nil {
		return errYoutubeClientCreate(err)
	}

	return nil
}
