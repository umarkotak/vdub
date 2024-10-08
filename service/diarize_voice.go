package service

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
)

func DiarizeVoice(ctx context.Context, taskDir string) error {
	cmd := exec.Command(
		"python", config.Get().PythonDiarizationPath,
		"--file_path", fmt.Sprintf("%s/raw_video_audio_Vocals_16KHz.wav", taskDir),
		"--output_path", fmt.Sprintf("%s/diarization.vtt", taskDir),
		"--auth_token", config.Get().HuggingFaceDiarizationToken,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"cmd": cmd.String(),
		}).Error(err)
		return err
	}
	return nil
}
