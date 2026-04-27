package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Credentials holds the API key and agent name
type Credentials struct {
	APIKey    string `json:"api_key"`
	AgentName string `json:"agent_name,omitempty"`
}

// GetConfigPath returns the path to the credentials file
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "moltbook", "credentials.json")
}

// LoadCredentials loads credentials from the config file
func LoadCredentials() (*Credentials, error) {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Credentials{}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &creds, nil
}

// SaveCredentials saves credentials to the config file
func SaveCredentials(creds *Credentials) error {
	path := GetConfigPath()
	dir := filepath.Dir(path)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetAPIKey returns the API key with priority: flag > env > config file
func GetAPIKey(flagKey string) (string, error) {
	// Priority 1: Command line flag
	if flagKey != "" {
		return flagKey, nil
	}

	// Priority 2: Environment variable
	if envKey := os.Getenv("MOLTBOOK_API_KEY"); envKey != "" {
		return envKey, nil
	}

	// Priority 3: Config file
	creds, err := LoadCredentials()
	if err != nil {
		return "", fmt.Errorf("failed to load credentials: %w", err)
	}

	if creds.APIKey == "" {
		return "", fmt.Errorf("no API key found. Set via --api-key flag, MOLTBOOK_API_KEY env var, or use 'molt auth set-key'")
	}

	return creds.APIKey, nil
}
