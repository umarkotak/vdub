package service

import (
	"context"
	"os"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/bregydoc/gtranslate"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

func TranslateTranscript(ctx context.Context, transcriptVttPath, targetResultPath string) error {
	vttContentByte, err := os.ReadFile(transcriptVttPath)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}
	vttContent := string(vttContentByte)

	subObj, err := astisub.OpenFile(transcriptVttPath)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	bar := progressbar.Default(int64(len(subObj.Items)), "Translating")
	for _, subItem := range subObj.Items {
		translated, err := gtranslate.TranslateWithParams(
			subItem.String(),
			gtranslate.TranslationParams{
				From: "en", To: "id",
			},
		)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}

		vttContent = strings.ReplaceAll(vttContent, subItem.String(), translated)
		bar.Add(1)
	}

	err = os.WriteFile(targetResultPath, []byte(vttContent), 0644)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}
