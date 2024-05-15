package service

import (
	"context"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/utils"
)

func SeparateVocal(ctx context.Context, rawAudioPath, targetDir string) error {
	cmd := exec.Command(
		"python", config.Get().VocalRemoverPy,
		"--input", rawAudioPath,
		"-P", config.Get().VocalRemoverModelPath,
		"-o", targetDir,
	)

	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"cmd": cmd.String(),
		}).Error(err)
		return err
	}

	utils.StreamStd(stderr)

	err = cmd.Wait()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"cmd": cmd.String(),
		}).Error(err)
		return err
	}

	return nil
}
