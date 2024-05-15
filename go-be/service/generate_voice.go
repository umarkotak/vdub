package service

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/asticode/go-astisub"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

const (
	VOICE_NAME  = "id-ID-ArdiNeural"
	VOICE_RATE  = "-10%"
	VOICE_PITCH = "-10Hz"
)

func GenerateVoice(ctx context.Context, transcriptTranslatedPath, targetSpeechDir string) error {
	cmd := exec.Command("mkdir", "-p", targetSpeechDir)
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	subObj, err := astisub.OpenFile(transcriptTranslatedPath)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	bar := progressbar.Default(int64(len(subObj.Items)), "Generating Audio")
	for idx, subItem := range subObj.Items {
		genSpeechPath := fmt.Sprintf("%s/%v.wav", targetSpeechDir, idx)
		cmd = exec.Command(
			"edge-tts",
			"--text", fmt.Sprintf("\"%s\"", subItem.String()),
			"--write-media", genSpeechPath,
			"-v", VOICE_NAME,
			fmt.Sprintf("--rate=%s", VOICE_RATE),
			fmt.Sprintf("--pitch=%s", VOICE_PITCH),
		)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		_, err = cmd.Output()
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"gen_speech_path": genSpeechPath,
				"cmd":             cmd.String(),
				"std_err":         stderr.String(),
			}).Error(err)
			return err
		}
		bar.Add(1)
	}

	return nil
}
