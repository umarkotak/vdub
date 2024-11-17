package service

import (
	"context"
	"os"

	"github.com/asticode/go-astisub"
	"github.com/bregydoc/gtranslate"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

func TranslateTranscript(ctx context.Context, taskDir, transcriptVttPath, targetResultPath string) error {
	logrusProc := logrus.WithContext(ctx).WithField("task_dir", taskDir)

	vttContentByte, err := os.ReadFile(transcriptVttPath)
	if err != nil {
		logrusProc.WithContext(ctx).Error(err)
		return err
	}
	vttContent := string(vttContentByte)

	subObj, err := astisub.OpenFile(transcriptVttPath)
	if err != nil {
		logrusProc.WithContext(ctx).Error(err)
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
			logrusProc.WithContext(ctx).Error(err)
			subItem.Lines[0].Items[0].Text = "error translate"
		} else {
			subItem.Lines[0].Items[0].Text = translated
		}

		bar.Add(1)
	}

	err = os.WriteFile(targetResultPath, []byte(vttContent), 0644)
	if err != nil {
		logrusProc.WithContext(ctx).Error(err)
		return err
	}

	return nil
}
