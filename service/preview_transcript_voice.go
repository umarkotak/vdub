package service

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/utils"
)

// Return full path binary of the preview voice
func GenPreviewTranscriptVoice(ctx context.Context, params model.TranscriptUpdatePosParams) (string, error) {
	taskDir := utils.GenTaskDir(params.TaskName)

	state, err := GetState(ctx, taskDir, model.TaskState{})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}

	translatedSub, err := astisub.OpenFile(utils.GenTranscriptTranslatedPath(taskDir))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}

	if params.Idx >= int64(len(translatedSub.Items)) {
		err = fmt.Errorf("invalid idx, out of bound")
		return "", err
	}

	genSpeechDir := fmt.Sprintf("%s/generated_speech", taskDir)
	adjustedSpeechDir := fmt.Sprintf("%s/adjusted_speech", taskDir)

	cmd := exec.Command("mkdir", "-p", genSpeechDir)
	_, err = cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}

	genSpeechPath := fmt.Sprintf("%s/%v.wav", genSpeechDir, params.Idx)
	subItem := translatedSub.Items[params.Idx]
	cmd = exec.Command(
		"/root/.pyenv/shims/edge-tts",
		"--text", fmt.Sprintf("\"%s\"", subItem.String()),
		"--write-media", genSpeechPath,
		"-v", state.VoiceName,
		fmt.Sprintf("--rate=%s", state.VoiceRate),
		fmt.Sprintf("--pitch=%s", state.VoicePitch),
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
		return "", err
	}

	originalDuration := subItem.EndAt - subItem.StartAt
	translatedDuration, _ := GetWavDuration(ctx, genSpeechPath)
	aTempo := utils.FloatToFixed(translatedDuration.Seconds()/originalDuration.Seconds(), 6)
	if aTempo < 1 {
		aTempo = 1.1
	} else if aTempo > 100 {
		aTempo = 100
	}

	adjustedSpeechPath := fmt.Sprintf("%s/%v.wav", adjustedSpeechDir, params.Idx)
	cmd = exec.Command(
		"ffmpeg", "-y",
		"-i", genSpeechPath,
		"-codec:a", "libmp3lame",
		"-filter:a", fmt.Sprintf("atempo=%v", aTempo),
		"-b:a", "320k",
		adjustedSpeechPath,
	)
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"cmd":     cmd.String(),
			"std_err": stderr.String(),
		}).Error(err)
		return "", err
	}

	return adjustedSpeechPath, nil
}

// Return full path binary of the preview voice
func GetPreviewTranscriptVoice(ctx context.Context, params model.TranscriptUpdatePosParams) (string, error) {
	taskDir := utils.GenTaskDir(params.TaskName)

	adjustedSpeechDir := fmt.Sprintf("%s/adjusted_speech", taskDir)

	adjustedSpeechPath := fmt.Sprintf("%s/%v.wav", adjustedSpeechDir, params.Idx)

	return adjustedSpeechPath, nil
}
