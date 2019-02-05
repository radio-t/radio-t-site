package youtube

import (
	"fmt"
	"os"
	"os/exec"
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
