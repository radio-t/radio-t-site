package youtube

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func makeVideo(audioPath, videoPath string) error {
	log.Info("Start making a video")

	f, err := ioutil.TempFile("", "intermedia_*.mp4")
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		if err := os.Remove(f.Name()); err != nil {
			log.Fatal(err)
		}
	}()

	createIntermediaVideoCommad := exec.Command("ffmpeg", "-y", "-loop", "1", "-i", "assets/cover.webp", "-c:v", "libx264", "-r", "15", "-pix_fmt", "yuv420p", "-t", "60", f.Name())

	if err := createIntermediaVideoCommad.Run(); err != nil {
		return errors.Wrap(err, "Error creating an intermedia file")
	}

	createVideoCommand := exec.Command("ffmpeg", "-y", "-stream_loop", "1", "-i", f.Name(), "-i", audioPath, "-c:v", "copy", "-c:a", "copy", "-shortest", videoPath)

	if err := createVideoCommand.Run(); err != nil {
		return errors.Wrap(err, "Error executing a ffmpeg")
	}

	log.Info("A video was made")
	return nil
}
