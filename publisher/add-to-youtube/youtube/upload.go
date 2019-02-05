package youtube

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/client"
	"golang.org/x/oauth2"
	yt "google.golang.org/api/youtube/v3"
)

func upload(config *oauth2.Config, audioPath, title, description, category, keywords, privacy, pathToSecrets string) (*yt.Video, error) {

	// prepare temprorary directory
	dir, err := ioutil.TempDir("", "add-to-youtube-")
	if err != nil {
		return nil, fmt.Errorf("Error creation a temprorary directory, got: %v", err)
	}
	defer os.RemoveAll(dir)

	baseName := path.Base(audioPath)
	videoPath := path.Join(dir, strings.TrimSuffix(baseName, filepath.Ext(baseName))+".mp4")
	if err := makeVideo(audioPath, videoPath); err != nil {
		return nil, err
	}

	// upload a video
	client, err := client.New(yt.YoutubeUploadScope, &client.Options{PathToSecrets: pathToSecrets, SkipAuth: true, Config: config})
	if err != nil {
		return nil, errYoutubeClientCreate(err)
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
		return nil, fmt.Errorf("Error opening %v: %v", videoPath, err)
	}

	// do an api request
	response, err := call.Media(file).Do()
	if err != nil {
		return nil, errAPICall(err)
	}

	return response, nil
}
