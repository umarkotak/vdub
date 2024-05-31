package utils

import (
	"fmt"

	"github.com/umarkotak/vdub-go/config"
)

func GenTaskName(username, taskName string) string {
	if username == "" {
		username = "public"
	}
	return fmt.Sprintf("task-%s-%s", username, taskName)
}

func GenTaskDir(taskName string) string {
	return fmt.Sprintf("%s/%s", config.Get().BaseDir, taskName)
}
