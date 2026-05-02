package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"molt/internal/config"
)

type LogEntry struct {
	Timestamp string   `json:"timestamp"`
	Input     []string `json:"input,omitempty"`
	Output    string   `json:"output"`
}

var (
	apiCalled   bool
	apiFailed   bool
	outputBuf   bytes.Buffer
	localConfig *config.LocalConfig
	mu          sync.Mutex
)

// Init initializes the logger with local configuration
func Init(cfg *config.LocalConfig) {
	localConfig = cfg
}

// SetAPICalled marks that an API call was attempted
func SetAPICalled() {
	mu.Lock()
	defer mu.Unlock()
	apiCalled = true
}

// SetAPIFailed marks that an API call failed
func SetAPIFailed() {
	mu.Lock()
	defer mu.Unlock()
	apiFailed = true
}

// CaptureOutput starts capturing Stdout and Stderr
func CaptureOutput() (func(), error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	stdout := os.Stdout
	stderr := os.Stderr

	os.Stdout = w
	os.Stderr = w

	done := make(chan bool)

	go func() {
		mw := io.MultiWriter(stdout, &outputBuf)
		_, _ = io.Copy(mw, r)
		done <- true
	}()

	return func() {
		_ = w.Close()
		<-done
		os.Stdout = stdout
		os.Stderr = stderr
	}, nil
}

// Finalize writes the logs based on the command outcome
func Finalize() error {
	mu.Lock()
	defer mu.Unlock()

	if localConfig == nil || !apiCalled {
		return nil
	}

	if localConfig.LogDir == "" {
		return nil
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(localConfig.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Output:    outputBuf.String(),
	}

	var filename string
	if apiFailed {
		if !localConfig.EnableErrorLog {
			return nil
		}
		filename = "error.log"
		entry.Input = os.Args
	} else {
		if !localConfig.EnableSuccessLog {
			return nil
		}
		filename = "log.log"
	}

	path := filepath.Join(localConfig.LogDir, filename)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}
