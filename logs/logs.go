package logs

import (
	"fmt"
	"os"
	"time"
)

type LogMessage struct {
	Time    time.Time
	Message string
}

var (
	errorLogChan  chan LogMessage
	accessLogChan chan LogMessage
)

// initialize starts the log writers and channels
func init() {
	errorLogChan = make(chan LogMessage, 100)
	accessLogChan = make(chan LogMessage, 100)

	go errorLogWriter(errorLogChan)
	go accessLogWriter(accessLogChan)
}

// LogError logs an error message
func LogError(err error) {
	if err != nil {
		errorLogChan <- LogMessage{
			Time:    time.Now(),
			Message: err.Error(),
		}
	}
}

// LogAccess logs an access message
func LogAccess(message string) {
	accessLogChan <- LogMessage{
		Time:    time.Now(),
		Message: message,
	}
}

func errorLogWriter(ch <-chan LogMessage) {
	filename := "storage/error_logs.txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open error log file:", err)
		return
	}
	defer file.Close()

	for msg := range ch {
		logLine := fmt.Sprintf("[%s] %s\n", msg.Time.Format(time.RFC3339), msg.Message)
		if _, err := file.WriteString(logLine); err != nil {
			fmt.Println("Failed to write to error log file:", err)
		}
	}
}

func accessLogWriter(ch <-chan LogMessage) {
	filename := "storage/access_logs.txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open access log file:", err)
		return
	}
	defer file.Close()

	for msg := range ch {
		logLine := fmt.Sprintf("[%s] %s\n", msg.Time.Format(time.RFC3339), msg.Message)
		if _, err := file.WriteString(logLine); err != nil {
			fmt.Println("Failed to write to access log file:", err)
		}
	}
}
