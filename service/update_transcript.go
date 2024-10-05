package service

import (
	"context"

	"github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/utils"
)

func UpdateTranscript(ctx context.Context, params model.TranscriptUpdateParams) error {
	taskDir := utils.GenTaskDir(params.TaskName)

	subObj, err := astisub.OpenFile(utils.GenTranscriptTranslatedPath(taskDir))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	for idx, subItem := range subObj.Items {
		subItem.StartAt, _ = utils.FormattedDurationToDuration(params.TranscriptData[idx].StartAt)
		subItem.EndAt, _ = utils.FormattedDurationToDuration(params.TranscriptData[idx].EndAt)
		subItem.Lines[0].Items[0].Text = params.TranscriptData[idx].Value
	}

	subObj.Write(utils.GenTranscriptTranslatedPath(taskDir))

	return nil
}
