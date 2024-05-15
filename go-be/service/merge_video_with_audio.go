package service

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func MergeVideoWithAudio(ctx context.Context, rawVideoPath, audioInstrumentPath, targetPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", rawVideoPath,
		"-i", audioInstrumentPath,
		"-c:v", "copy",
		"-c:a", "aac",
		"-map", "0:v:0",
		"-map", "1:a:0",
		targetPath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"raw_video_path":        rawVideoPath,
			"audio_instrument_path": audioInstrumentPath,
			"std_err":               stderr.String(),
		}).Error(err)
		return err
	}
	return nil
}
