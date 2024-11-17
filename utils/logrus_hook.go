package utils

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type CustomLogrusHook struct {
	LogFilePath string
}

// Levels defines the log levels that this hook will process
func (hook *CustomLogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels // Log all levels
}

// Fire is called when a log event occurs
func (hook *CustomLogrusHook) Fire(entry *logrus.Entry) error {
	taskDir := fmt.Sprint(entry.Data["task_dir"])

	if taskDir == "" {
		return nil
	}

	QuickStoreLog(taskDir, strings.ToUpper(entry.Level.String()), fmt.Sprintf("%s: %+v", entry.Message, entry.Data))

	return nil
}
