package service

import (
	"context"
	"fmt"

	"github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/utils"
)

type (
	TranscriptInfo struct {
		TranscriptLines []TranscriptLine `json:"transcript_line"`
	}

	TranscriptLine struct {
		StartAt string `json:"start_at"`
		EndAt   string `json:"end_at"`
		Value   string `json:"value"`
	}
)

func GetTranscript(ctx context.Context, taskName, transcriptType string) (TranscriptInfo, error) {
	taskDir := utils.GenTaskDir(taskName)
	transcriptFileName := "transcript.vtt"
	if transcriptType == "translated" {
		transcriptFileName = "transcript_translated.vtt"
	}
	transcriptPath := fmt.Sprintf("%s/%s", taskDir, transcriptFileName)

	transcriptSubtitle, err := astisub.OpenFile(transcriptPath)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return TranscriptInfo{}, nil
	}

	transcriptLines := []TranscriptLine{}
	for _, oneLine := range transcriptSubtitle.Items {
		transcriptLines = append(transcriptLines, TranscriptLine{
			StartAt: utils.DurationToFormattedDuration(oneLine.StartAt),
			EndAt:   utils.DurationToFormattedDuration(oneLine.EndAt),
			Value:   oneLine.String(),
		})
	}

	return TranscriptInfo{
		TranscriptLines: transcriptLines,
	}, nil
}
