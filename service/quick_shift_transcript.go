package service

import (
	"context"
	"time"

	"github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/utils"
)

const (
	MIN_EXPECTED_CHAR_PER_SEC = float64(17.0)
	MAX_EXPECTED_CHAR_PER_SEC = float64(18.0)
)

func QuickShiftTranscript(ctx context.Context, params model.TranscriptUpdateParams) error {
	taskDir := utils.GenTaskDir(params.TaskName)

	subObj, err := astisub.OpenFile(utils.GenTranscriptTranslatedPath(taskDir))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	for idx, subItem := range subObj.Items {
		charLength := float64(len(subItem.String()))
		duration := (subItem.EndAt - subItem.StartAt)
		charPerSec := charLength / duration.Seconds()

		if idx > 0 {
			prevSubItem := subObj.Items[idx-1]

			if prevSubItem.EndAt > subItem.StartAt {
				subItem.StartAt = prevSubItem.EndAt
			}

			if subItem.StartAt > prevSubItem.EndAt {
				if (subItem.StartAt - prevSubItem.EndAt) <= 3*time.Second {
					subItem.StartAt = prevSubItem.EndAt
				}
			}

			duration = (subItem.EndAt - subItem.StartAt)
			charPerSec = charLength / duration.Seconds()
		}

		if charPerSec > MAX_EXPECTED_CHAR_PER_SEC {
			subItem.EndAt = subItem.StartAt + time.Duration((charLength/MIN_EXPECTED_CHAR_PER_SEC)*float64(time.Second))
		} else if charPerSec < MIN_EXPECTED_CHAR_PER_SEC {
			subItem.EndAt = subItem.StartAt + time.Duration((charLength/MIN_EXPECTED_CHAR_PER_SEC)*float64(time.Second))
		}

		subObj.Items[idx] = subItem
	}

	subObj.Write(utils.GenTranscriptTranslatedPath(taskDir))

	return nil
}
