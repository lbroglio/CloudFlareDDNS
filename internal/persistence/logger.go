package persistence

import (
	"fmt"
	"os"
	"time"

	"github.com/lbroglio/CloudFlareDDNS/internal/config"
)

// ensureLogDirectoryExists checks if the log directory exists, and creates it if it doesn't.
func ensureLogDirectoryExists() error {
	logDirectory := config.GetPersistenceRootPath() + "/logs"
	info, err := os.Stat(logDirectory)

	if err != nil {
		os.MkdirAll(logDirectory, os.ModePerm)
	} else {
		if !info.IsDir() {
			return fmt.Errorf("log path exists but is not a directory: %s", logDirectory)
		}
	}

	return nil
}

// getLogFilePath returns the path to the log file for the current date.
func getLogFilePath() string {
	timestamp := time.Now().Format("2006-01-02")
	return config.GetPersistenceRootPath() + fmt.Sprintf("/logs/%s.log", timestamp)
}

// writeLogLine appends a log line to the log file for the current date.
func WriteLogLine(contents string) error {
	err := ensureLogDirectoryExists()

	if err != nil {
		return fmt.Errorf("failed to ensure log directory exists: %v", err)
	}

	logFilePath := getLogFilePath()
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	stringToWrite := time.Now().Format("2006-01-02 15:04:05") + " - " + contents
	file.WriteString(stringToWrite + "\n")

	return nil

}
