package service

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
)

func DeleteTask(ctx context.Context, taskName string) error {
	if taskName == "" {
		return fmt.Errorf("missing task name")
	}
	cmd := exec.Command(
		"rm", "-rf", fmt.Sprintf("%s/%s", config.Get().BaseDir, taskName),
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}
	return nil
}
