package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

// ConfigManager handles loading and managing configurations.
type ConfigManager struct {
	config map[string]interface{}
}

// NewConfigManager creates a new instance of ConfigManager.
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		config: make(map[string]interface{}),
	}
}

// LoadConfig loads configuration from a JSON file.
func (cm *ConfigManager) LoadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cm.config); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}

	return nil
}

// Get retrieves a configuration value by key.
func (cm *ConfigManager) Get(key string) (interface{}, error) {
	value, exists := cm.config[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found in configuration", key)
	}
	return value, nil
}

// GetString retrieves a string configuration value by key.
func (cm *ConfigManager) GetString(key string) (string, error) {
	value, err := cm.Get(key)
	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string", key)
	}
	return strValue, nil
}