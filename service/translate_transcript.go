package service

import (
	"context"

	"github.com/asticode/go-astisub"
	"github.com/bregydoc/gtranslate"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

func TranslateTranscript(ctx context.Context, taskDir, transcriptVttPath, targetResultPath string) error {
	logrusProc := logrus.WithContext(ctx).WithField("task_dir", taskDir)

	subObj, err := astisub.OpenFile(transcriptVttPath)
	if err != nil {
		logrusProc.WithContext(ctx).Error(err)
		return err
	}

	bar := progressbar.Default(int64(len(subObj.Items)), "Translating")
	for _, subItem := range subObj.Items {
		// originalText := subItem.String()
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

		// logrusProc.WithContext(ctx).Infof("translate: \"%s\" into \"%s\"", originalText, translated)

		bar.Add(1)
	}

	subObj.Write(targetResultPath)

	return nil
}
