package service

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/asticode/go-astisub"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/utils"
)

func AdjustVoiceSpeed(ctx context.Context, transcriptPath, originalSpeechDir, adjustedSpeechDir string) error {
	cmd := exec.Command("mkdir", "-p", adjustedSpeechDir)
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	subObj, err := astisub.OpenFile(transcriptPath)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	bar := progressbar.Default(int64(len(subObj.Items)), "Adjusting Audio")
	for idx, subItem := range subObj.Items {
		genSpeechPath := fmt.Sprintf("%s/%v.wav", originalSpeechDir, idx)
		adjustedSpeechPath := fmt.Sprintf("%s/%v.wav", adjustedSpeechDir, idx)

		originalDuration := subItem.EndAt - subItem.StartAt
		translatedDuration, _ := GetWavDuration(ctx, genSpeechPath)
		aTempo := utils.FloatToFixed(translatedDuration.Seconds()/originalDuration.Seconds(), 6)
		if aTempo < 1 {
			aTempo = 1.1
		} else if aTempo > 100 {
			aTempo = 100
		}

		cmd = exec.Command(
			"ffmpeg",
			"-i", genSpeechPath,
			"-codec:a", "libmp3lame",
			"-filter:a", fmt.Sprintf("atempo=%v", aTempo),
			"-b:a", "320k",
			adjustedSpeechPath,
		)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		_, err = cmd.Output()
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"std_err": stderr.String(),
			}).Error(err)
			return err
		}
		bar.Add(1)
	}

	return nil
}

func GetWavDuration(ctx context.Context, filename string) (time.Duration, error) {
	cmd := exec.Command(
		"ffprobe",
		"-i", filename,
		"-show_entries",
		"format=duration",
		"-v", "quiet",
		"-of",
		"csv=p=0",
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"std_err": stderr.String(),
		}).Error(err)
		return 0, err
	}

	seconds, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return 0, err
	}

	// Convert seconds to time.Duration (in nanoseconds)
	duration := time.Duration(seconds * float64(time.Second))

	return duration, nil
}
