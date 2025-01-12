package service

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
)

func MergeVideoWithDubb(
	ctx context.Context,
	transcriptTranslatedPath,
	adjustedSpeechDir,
	instrumentVideoPath,
	dubbedVideoPath string,
	volume string,
) error {
	subObj, err := astisub.OpenFile(transcriptTranslatedPath)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	ffmpegArgs := []string{
		"-y", "-i", instrumentVideoPath,
	}

	for idx := range subObj.Items {
		adjustedSpeechPath := fmt.Sprintf("%s/%v.wav", adjustedSpeechDir, idx)
		ffmpegArgs = append(ffmpegArgs, "-i", adjustedSpeechPath)
	}

	filterComplexes := []string{
		"[0]volume=10dB[video0]",
	}
	filterComplexCloser := ""
	for idx, subItem := range subObj.Items {
		audioIdx := fmt.Sprintf("[audio%v]", idx)

		filter := fmt.Sprintf(
			"[%v]volume=10dB,adelay=%v%s",
			idx+1,
			subItem.StartAt.Milliseconds(),
			audioIdx,
		)
		filterComplexes = append(filterComplexes, filter)

		filterComplexCloser += audioIdx
	}
	filterComplexCloserFormatted := fmt.Sprintf("[video0]%samix=%v", filterComplexCloser, len(subObj.Items)+1)
	filterComplexes = append(filterComplexes, filterComplexCloserFormatted)

	ffmpegArgs = append(ffmpegArgs, "-filter_complex", strings.Join(filterComplexes, ";"))

	ffmpegArgs = append(
		ffmpegArgs,
		"-c:v", "copy",
		dubbedVideoPath,
	)

	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd":     cmd.String(),
			"std_err": stderr.String(),
		}).Error(err)
		return err
	}

	// used to increase video volume
	// ffmpeg -i input.mp4 -filter:a "volume=2.0" -c:v copy output.mp4
	if volume == "" {
		volume = "1.2"
	}
	cmd = exec.Command("ffmpeg",
		"-i", dubbedVideoPath,
		"-filter:a", fmt.Sprintf("\"volume=%s\"", volume),
		"-c:v", "copy", dubbedVideoPath,
	)
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd":     cmd.String(),
			"std_err": stderr.String(),
		}).Error(err)
		return err
	}

	return nil
}
