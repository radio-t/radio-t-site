package youtube

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/client"
	yt "google.golang.org/api/youtube/v3"
)

func upload(c *Client, audioPath, title, description, category, keywords, privacy string) (*yt.Video, error) {

	log.Info("Creating temporary directory")
	dir, err := ioutil.TempDir("", "add-to-youtube-")
	if err != nil {
		return nil, errors.Wrap(err, "Error creation a temprorary directory")
	}
	defer func() {
		log.Infof("Removing temporary directory `%s`", dir)
		os.RemoveAll(dir)
	}()

	baseName := path.Base(audioPath)
	videoPath := path.Join(dir, strings.TrimSuffix(baseName, filepath.Ext(baseName))+".mp4")
	if err := makeVideo(audioPath, c.CoverPath, videoPath); err != nil {
		return nil, err
	}

	// upload a video
	client, err := client.New(c.OAuth2, c.TokenPath, true, c.Scopes...)
	if err != nil {
		return nil, errOAuth2HTTPClientCreate(err)
	}

	// create youtube client
	service, err := yt.New(client)
	if err != nil {
		return nil, errYoutubeClientCreate(err)
	}

	// prepare a metas for video
	upload := &yt.Video{
		Snippet: &yt.VideoSnippet{
			Title:                title,
			Description:          description,
			CategoryId:           category,
			DefaultAudioLanguage: "ru",
			DefaultLanguage:      "ru",
		},
		Status: &yt.VideoStatus{PrivacyStatus: privacy, License: "creativeCommon"},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}

	// prepare an api request
	call := service.Videos.Insert("snippet,status", upload)
	file, err := os.Open(videoPath)
	defer file.Close()
	if err != nil {
		return nil, errors.Errorf("Error opening %v: %v", videoPath, err)
	}

	// do an api request
	log.Info("Uploading video to YouTube")
	response, err := call.Media(file).Do()
	if err != nil {
		return nil, errors.Wrap(err, "Error making an API call")
	}

	return response, nil
}
