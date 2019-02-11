package youtube

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func makeVideo(audioPath, videoPath string) error {
	log.Info("Start making a video")

	exec := exec.Command("ffmpeg", "-loop", "1", "-i", "assets/cover.webp", "-i", audioPath, "-c:v", "libx264", "-r", "15", "-c:a", "copy", "-shortest", "-y", "-pix_fmt", "yuv420p", videoPath)

	exec.Stdout = os.Stdout
	exec.Stderr = os.Stderr

	if err := exec.Run(); err != nil {
		return errors.Wrap(err, "Error executing a ffmpeg")
	}

	log.Info("A video was made")
	return nil
}
