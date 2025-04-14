package logs

import (
	"fmt"
	"os"
	"time"
)

type LogType string

type LogMessage struct {
	Time    time.Time
	Message string
}

func ErrorLogWriter(errorLogChan <-chan LogMessage) {
	filename := "storage/error_logs.txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()

	for msg := range errorLogChan {
		logLine := fmt.Sprintf("[%s] %s\n", msg.Time.Format(time.RFC3339), msg.Message)
		_, err := file.WriteString(logLine)
		if err != nil {
			fmt.Println("Failed to write to log file:", err)
		}
	}
}

func AccessLogWriter(accessLogChan <-chan LogMessage) {
	filename := "storage/access_logs.txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()

	for msg := range accessLogChan {
		logLine := fmt.Sprintf("[%s] %s\n", msg.Time.Format(time.RFC3339), msg.Message)
		_, err := file.WriteString(logLine)
		if err != nil {
			fmt.Println("Failed to write to log file:", err)
		}
	}
}

func CreateErrorLogMessage(err error) LogMessage {
	return LogMessage{
		Time:    time.Now(),
		Message: err.Error(),
	}
}

func CreateAccessLogMessage(message string) LogMessage {
	return LogMessage{
		Time:    time.Now(),
		Message: message,
	}
}
