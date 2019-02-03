package youtube

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/client"
	"google.golang.org/api/youtube/v3"
)

func makeVideo(audioPath, videoPath string) error {
	exec := exec.Command("ffmpeg", "-loop", "1", "-i", "assets/cover.webp", "-i", audioPath, "-c:v", "libx264", "-r", "15", "-c:a", "copy", "-shortest", "-y", "-pix_fmt", "yuv420p", videoPath)
	exec.Stdout = os.Stdout
	exec.Stderr = os.Stderr
	if err := exec.Run(); err != nil {
		return fmt.Errorf("Error making a video, got: %v", err)
	}
	return nil
}

// Upload uses an audio file to create a video file, then upload it with metadatas to Youtube.
func Upload(audioPath, title, description, category, keywords, privacy, pathToSecrets string) (*youtube.Video, error) {

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
	client, err := client.New(youtube.YoutubeUploadScope, &client.Options{PathToSecrets: pathToSecrets, SkipAuth: true})
	if err != nil {
		return nil, errYoutubeClientCreate(err)
	}

	// create youtube client
	service, err := youtube.New(client)
	if err != nil {
		return nil, errYoutubeClientCreate(err)
	}

	// prepare a metas for video
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:                title,
			Description:          description,
			CategoryId:           category,
			DefaultAudioLanguage: "ru",
			DefaultLanguage:      "ru",
		},
		Status: &youtube.VideoStatus{PrivacyStatus: privacy, License: "creativeCommon"},
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
