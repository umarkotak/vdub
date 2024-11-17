package service

import (
	"context"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/utils"
)

// Separate vocal and audio
func SeparateVocal(ctx context.Context, rawAudioPath, targetDir string) error {
	cmd := exec.Command(
		"audio-separator",
		"--model_filename", "UVR-MDX-NET-Voc_FT.onnx",
		"--output_format", "wav",
		"--output_dir", targetDir,
		rawAudioPath,
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
