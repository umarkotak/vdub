package service

import (
	"context"

	"github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/utils"
)

func AddTranscriptByIdx(ctx context.Context, params model.TranscriptUpdatePosParams) error {
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

	originalSub.Items = addItemAfterIndex(originalSub.Items, originalSub.Items[params.Idx], params.Idx)
	originalSub.Write(utils.GenTranscriptVttPath(taskDir))

	translatedSub.Items = addItemAfterIndex(translatedSub.Items, translatedSub.Items[params.Idx], params.Idx)
	translatedSub.Write(utils.GenTranscriptTranslatedPath(taskDir))

	return nil
}

func addItemAfterIndex(slice []*astisub.Item, newItem *astisub.Item, index int64) []*astisub.Item {
	return append(slice[:index+1], append([]*astisub.Item{newItem}, slice[index+1:]...)...)
}
