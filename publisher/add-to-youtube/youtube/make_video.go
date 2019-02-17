package youtube

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func makeVideo(audioPath, coverPath, videoPath string) error {
	log.Info("Making a video")

	dir, err := ioutil.TempDir("", "make_intermedia_video_")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	f := path.Join(dir, "intermedia.mp4")

	createIntermediaVideoCommad := exec.Command("ffmpeg", "-y", "-loop", "1", "-i", coverPath, "-c:v", "libx264", "-r", "15", "-pix_fmt", "yuv420p", "-t", "60", f)
	// createIntermediaVideoCommad.Stderr = os.Stderr
	// createIntermediaVideoCommad.Stdout = os.Stdout

	if err := createIntermediaVideoCommad.Run(); err != nil {
		return errors.Wrap(err, "Error creating an intermedia video file")
	}

	createVideoCommand := exec.Command("ffmpeg", "-y", "-stream_loop", "-1", "-i", f, "-i", audioPath, "-c", "copy", "-shortest", videoPath)

	if err := createVideoCommand.Run(); err != nil {
		return errors.Wrap(err, "Error creating a podcast episode video")
	}
	// createVideoCommand.Stderr = os.Stderr
	// createVideoCommand.Stdout = os.Stdout

	log.Infof("A video was made at `%s`", videoPath)
	return nil
}
