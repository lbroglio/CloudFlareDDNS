package persistence

import (
	"os"
	"testing"
	"time"

	"github.com/lbroglio/CloudFlareDDNS/internal/config"
)

func loggerTestSetup(t *testing.T) {
	testPeristenceRoot, _ := os.MkdirTemp(os.TempDir(), "test_persistence_root")
	t.Setenv("CLOUDFLAREDDNS_PERSISTENCE_ROOT", testPeristenceRoot)
}

func TestEnsureLogDirectoryExists(t *testing.T) {
	loggerTestSetup(t)

	err := ensureLogDirectoryExists()
	if err != nil {
		t.Errorf("ensureLogDirectoryExists() returned an error: %v", err)
	}

	logDirectory := config.GetPersistenceRootPath() + "/logs"
	info, err := os.Stat(logDirectory)
	if err != nil {
		t.Errorf("Expected %s to exist: %v", logDirectory, err)
	}
	if !info.IsDir() {
		t.Errorf("Expected %s to be a directory", logDirectory)
	}
}

func TestGetLogFilePath(t *testing.T) {
	loggerTestSetup(t)

	logFilePath := getLogFilePath()
	date := time.Now().Format("2006-01-02")
	expectedPath := config.GetPersistenceRootPath() + "/logs/" + date + ".log"

	if logFilePath != expectedPath {
		t.Errorf("Expected log file path to be %s, but got %s", expectedPath, logFilePath)
	}

}

func TestWriteLogLine(t *testing.T) {
	loggerTestSetup(t)

	logLine := "This is a test log line."
	err := WriteLogLine(logLine)

	if err != nil {
		t.Errorf("writeLogLine() returned an error: %v", err)
	}

	logFilePath := getLogFilePath()
	bytes, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Errorf("Failed to read log file: %v", err)
	}

	contents := string(bytes)
	expectedLine := time.Now().Format("2006-01-02 15:04:05") + " - " + logLine + "\n"
	if contents != expectedLine {
		t.Errorf("Expected log line to be '%s', but got '%s'", expectedLine, contents)
	}

}
