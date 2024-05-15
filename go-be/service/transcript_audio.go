package service

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/utils"
)

const (
	BEAM_SIZE      = "6"
	ENTHROPY_THOLD = "2.8"
	MAX_CONTEXT    = "128"
)

func TranscriptAudio(ctx context.Context, audioPath, transcriptPath string) error {
	cmdTranscript := exec.Command(
		config.Get().WhisperBinary,
		"-m", config.Get().WhisperModelPath,
		"-ovtt",
		"--beam-size", BEAM_SIZE,
		"--entropy-thold", ENTHROPY_THOLD,
		"--max-context", MAX_CONTEXT,
		"-of", transcriptPath,
		audioPath,
	)

	stdout, _ := cmdTranscript.StdoutPipe()
	stderr, _ := cmdTranscript.StderrPipe()
	err := cmdTranscript.Start()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"cmd": cmdTranscript.String(),
		}).Error(err)
	}

	utils.StreamCmdTranscript(stdout, stderr)

	err = cmdTranscript.Wait()
	fmt.Printf("\n")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd": cmdTranscript.String(),
		}).Error(err)
	}

	return nil
}
