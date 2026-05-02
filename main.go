package main

import (
	"fmt"
	"os"

	"molt/cmd"
	"molt/internal/config"
	"molt/internal/logger"
)

func main() {
	// Load local config
	localConfig, err := config.LoadLocalConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading local config: %v\n", err)
	}

	// Initialize logger
	logger.Init(localConfig)

	// Start output capture
	stopCapture, err := logger.CaptureOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting output capture: %v\n", err)
	}

	// Execute command
	execErr := cmd.Execute()

	// Stop capture and finalize logs
	if stopCapture != nil {
		stopCapture()
	}

	if err := logger.Finalize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error finalizing logs: %v\n", err)
	}

	if execErr != nil {
		os.Exit(1)
	}
}
