package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

func QuickStoreLog(taskDir, logType, logContent string) error {
	return StoreLog(
		fmt.Sprintf("%s/log.log", taskDir),
		logType,
		logContent,
	)
}

func StoreLog(logFilePath, logType, logContent string) error {
	// 1. Open the log file in append mode:
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // 0644: Permissions for the file
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	defer file.Close()

	// 2. Get the current timestamp:
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 3. Format the log message:
	logMessage := fmt.Sprintf("%s | %s | %s\n", timestamp, logType, logContent)

	// 4. Append the log message to the file:
	if _, err := file.WriteString(logMessage); err != nil {
		return fmt.Errorf("error writing to log file: %w", err)
	}

	return nil
}

func QuickGetLog(taskDir string) ([]LogEntry, error) {
	logFilePath := fmt.Sprintf("%s/log.log", taskDir)

	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}
	defer file.Close()

	var logEntries []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " | ", 3) // Split into at most 3 parts

		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid log format: %s", line)
		}

		timestamp, err := time.Parse("2006-01-02 15:04:05", parts[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp: %w", err)
		}

		entry := LogEntry{
			Timestamp: timestamp,
			Level:     parts[1],
			Message:   parts[2],
		}
		logEntries = append(logEntries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	return logEntries, nil
}
