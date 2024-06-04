package service

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func GenerateVideoSnapshot(ctx context.Context, rawVideoPath, targetPath string) error {
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-ss", "00:00:01",
		"-i", rawVideoPath,
		"-frames:v", "1",
		"-q:v", "2",
		targetPath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"raw_video_path":       rawVideoPath,
			"target_snapshot_path": targetPath,
			"std_err":              stderr.String(),
		}).Error(err)
		return err
	}
	return nil
}
