package service

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func DownloadYoutubeVideo(ctx context.Context, videoUrl, targetPath string) error {
	cmd := exec.Command(
		"yt-dlp",
		"--progress",
		"-S", "ext",
		"--cookies", "cookies.txt",
		"-o", targetPath,
		videoUrl,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"std_err": stderr.String(),
		}).Error(err)
		return err
	}
	return nil
}
