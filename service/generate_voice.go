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

type (
	VoiceOpts struct {
		Name  string
		Rate  string
		Pitch string
	}
)

func GenerateVoice(ctx context.Context, taskDir, transcriptTranslatedPath, targetSpeechDir string, voiceOpts VoiceOpts) error {
	logrusProc := logrus.WithContext(ctx).WithField("task_dir", taskDir)

	voiceOpts.SetDefault()

	cmd := exec.Command("mkdir", "-p", targetSpeechDir)
	_, err := cmd.Output()
	if err != nil {
		logrusProc.WithContext(ctx).Error(err)
		return err
	}

	subObj, err := astisub.OpenFile(transcriptTranslatedPath)
	if err != nil {
		logrusProc.WithContext(ctx).Error(err)
		return err
	}

	bar := progressbar.Default(int64(len(subObj.Items)), "Generating Audio")
	for idx, subItem := range subObj.Items {
		genSpeechPath := fmt.Sprintf("%s/%v.wav", targetSpeechDir, idx)
		cmd = exec.Command(
			"/root/.pyenv/shims/edge-tts",
			"--text", fmt.Sprintf("\"%s\"", subItem.String()),
			"--write-media", genSpeechPath,
			"-v", voiceOpts.Name,
			fmt.Sprintf("--rate=%s", voiceOpts.Rate),
			fmt.Sprintf("--pitch=%s", voiceOpts.Pitch),
		)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		_, err = cmd.Output()
		if err != nil {
			logrusProc.WithContext(ctx).WithFields(logrus.Fields{
				"gen_speech_path": genSpeechPath,
				"cmd":             cmd.String(),
				"std_err":         stderr.String(),
			}).Error(err)

			// fallback when generating voice
			cmd = exec.Command(
				"/root/.pyenv/shims/edge-tts",
				"--text", fmt.Sprintf("\"%s\"", "error tts"),
				"--write-media", genSpeechPath,
				"-v", voiceOpts.Name,
				fmt.Sprintf("--rate=%s", voiceOpts.Rate),
				fmt.Sprintf("--pitch=%s", voiceOpts.Pitch),
			)
			_, err = cmd.Output()
			if err != nil {
				logrusProc.WithContext(ctx).WithFields(logrus.Fields{
					"gen_speech_path": genSpeechPath,
					"cmd":             cmd.String(),
					"std_err":         stderr.String(),
					"fallback_err":    true,
				}).Error(err)
				return err
			}
		}
		bar.Add(1)
	}

	return nil
}

func (vo *VoiceOpts) SetDefault() {
	if vo.Name == "" {
		vo.Name = VOICE_NAME
	}

	if vo.Rate == "" {
		vo.Rate = VOICE_RATE
	}

	if vo.Pitch == "" {
		vo.Pitch = VOICE_PITCH
	}
}
