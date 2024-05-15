package service

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func Generate16KHzAudio(ctx context.Context, audioPath, targetPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", audioPath,
		"-acodec", "pcm_s16le",
		"-ac", "1",
		"-ar", "16000",
		targetPath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"audio_path": audioPath,
			"std_err":    stderr.String(),
		}).Error(err)
		return err
	}

	return nil
}
