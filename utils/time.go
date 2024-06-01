package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func DurationToFormattedDuration(d time.Duration) string {
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	milliseconds := d % time.Second / time.Millisecond

	return fmt.Sprintf("%02d:%02d:%02d:%03d", hours, minutes, seconds, milliseconds)
}

func FormattedDurationToDuration(d string) (time.Duration, error) {
	parts := strings.Split(d, ":")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid duration format: %s", d)
	}

	hours, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid hours: %s", parts[0])
	}

	minutes, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %s", parts[1])
	}

	seconds, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %s", parts[2])
	}

	milliseconds, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid milliseconds: %s", parts[3])
	}

	// Calculate the total duration in nanoseconds
	duration := (time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(milliseconds)*time.Millisecond).Nanoseconds()

	return time.Duration(duration), nil
}
