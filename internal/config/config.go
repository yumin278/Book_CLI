package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LocalConfig holds configuration from the executable's directory
type LocalConfig struct {
	APIKey           string `json:"api_key"`
	Proxy            string `json:"proxy,omitempty"`
	LogDir           string `json:"log_dir"`
	EnableSuccessLog bool   `json:"enable_success_log"`
	EnableErrorLog   bool   `json:"enable_error_log"`
	CognitiveRoot    string `json:"cognitive_root,omitempty"`
	LLMEndpoint      string `json:"llm_endpoint,omitempty"`
	LLMModel         string `json:"llm_model,omitempty"`
	LLMAPIKey        string `json:"llm_api_key,omitempty"`
}

// Credentials holds the API key and agent name
type Credentials struct {
	APIKey string `json:"api_key"`

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

// LoadLocalConfig loads configuration from config.json in the executable's directory
func LoadLocalConfig() (*LocalConfig, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	path := filepath.Join(filepath.Dir(exePath), "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Not found is acceptable
		}
		return nil, fmt.Errorf("failed to read local config: %w", err)
	}

	var config LocalConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse local config: %w", err)
	}

	return &config, nil
}

// GetAPIKey returns the API key with priority: local config > flag > env > global config file
func GetAPIKey(flagKey string) (string, error) {
	// Priority 0: Local config in program directory
	localConfig, _ := LoadLocalConfig()
	if localConfig != nil && localConfig.APIKey != "" {
		return localConfig.APIKey, nil
	}

	// Priority 1: Command line flag
	if flagKey != "" {
		return flagKey, nil
	}

	// Priority 2: Environment variable
	if envKey := os.Getenv("MOLTBOOK_API_KEY"); envKey != "" {
		return envKey, nil
	}

	// Priority 3: Global config file
	creds, err := LoadCredentials()
	if err != nil {
		return "", fmt.Errorf("failed to load credentials: %w", err)
	}

	if creds.APIKey == "" {
		return "", fmt.Errorf("no API key found. Set via --api-key flag, MOLTBOOK_API_KEY env var, or use 'molt auth set-key'")
	}

	return creds.APIKey, nil
}

// GetProxyURL returns the proxy URL with priority: flag > env > local config > global config file
func GetProxyURL(flagProxy string) string {
	// Priority 1: Command line flag
	if flagProxy != "" {
		return flagProxy
	}

	// Priority 2: Environment variable
	if envProxy := os.Getenv("MOLTBOOK_PROXY"); envProxy != "" {
		return envProxy
	}

	// Priority 3: Local config in program directory
	localConfig, _ := LoadLocalConfig()
	if localConfig != nil && localConfig.Proxy != "" {
		return localConfig.Proxy
	}

	return ""
}
