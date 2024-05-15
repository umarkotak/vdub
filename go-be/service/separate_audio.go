package service

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func SeparateAudio(ctx context.Context, rawAudioPath, targetAudioPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", rawAudioPath,
		"-vn",
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2",
		targetAudioPath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"raw_audio_path":    rawAudioPath,
			"target_audio_path": targetAudioPath,
			"std_err":           stderr.String(),
		}).Error(err)
		return err
	}
	return nil
}
