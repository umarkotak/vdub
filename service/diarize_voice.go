package service

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/utils"
)

func DiarizeVoice(ctx context.Context, taskDir string) error {
	cmd := exec.Command(
		"python", config.Get().PythonDiarizationPath,
		"--file_path", utils.GenVocal16KHzPath(taskDir),
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
