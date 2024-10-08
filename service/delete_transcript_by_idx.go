package service

import (
	"context"

	"github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/utils"
)

func DeleteTranscriptByIdx(ctx context.Context, params model.TranscriptUpdatePosParams) error {
	taskDir := utils.GenTaskDir(params.TaskName)

	originalSub, err := astisub.OpenFile(utils.GenTranscriptVttPath(taskDir))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	translatedSub, err := astisub.OpenFile(utils.GenTranscriptTranslatedPath(taskDir))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	originalSub.Items = removeSubItemByIndex(originalSub.Items, params.Idx)
	originalSub.Write(utils.GenTranscriptVttPath(taskDir))

	translatedSub.Items = removeSubItemByIndex(translatedSub.Items, params.Idx)
	translatedSub.Write(utils.GenTranscriptTranslatedPath(taskDir))

	return nil
}

func removeSubItemByIndex(slice []*astisub.Item, index int64) []*astisub.Item {
	return append(slice[:index], slice[index+1:]...)
}
