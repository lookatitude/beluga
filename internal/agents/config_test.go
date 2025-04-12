package agents

import (
	"beluga/pkg/agents/config"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigManager(t *testing.T) {
	// Create a temporary test config file
	testConfig := createTestConfig()
	tempDir, err := ioutil.TempDir("", "beluga-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // clean up after test
	
	configPath := filepath.Join(tempDir, "agents.json")
	if err := writeTestConfig(configPath, testConfig); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	
	// Create and initialize config manager
	cm := config.NewConfigManager()
	if err := cm.LoadConfig(configPath); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Test GetAgentConfig
	testCases := []struct {
		name     string
		agentName string
		shouldExist bool
	}{
		{"existing agent", "test_agent", true},
		{"non-existing agent", "unknown_agent", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			agentConfig, err := cm.GetAgentConfig(tc.agentName)
			if tc.shouldExist {
				if err != nil {
					t.Errorf("Expected agent %s to exist, got error: %v", tc.agentName, err)
				}
				if agentConfig == nil {
					t.Errorf("Agent config should not be nil for %s", tc.agentName)
				} else if agentConfig.Name != tc.agentName {
					t.Errorf("Expected agent name %s, got %s", tc.agentName, agentConfig.Name)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for non-existent agent %s, got nil", tc.agentName)
				}
			}
		})
	}
	
	// Test GetAllAgentConfigs
	allConfigs := cm.GetAllAgentConfigs()
	if len(allConfigs) != 2 {
		t.Errorf("Expected 2 agent configs, got %d", len(allConfigs))
	}
	
	// Test default settings
	defaultSettings := cm.GetDefaultSettings()
	if defaultSettings["log_level"] != "debug" {
		t.Errorf("Expected default log_level to be 'debug', got '%v'", defaultSettings["log_level"])
	}
	
	// Test logging config
	loggingConfig := cm.GetLoggingConfig()
	if loggingConfig["console_enabled"] != true {
		t.Errorf("Expected console_enabled to be true, got %v", loggingConfig["console_enabled"])
	}
	
	// Test SaveConfig
	newConfigPath := filepath.Join(tempDir, "updated_config.json")
	if err := cm.SaveConfig(newConfigPath); err != nil {
		t.Errorf("Failed to save config: %v", err)
	}
	
	// Verify saved file exists
	if _, err := os.Stat(newConfigPath); os.IsNotExist(err) {
		t.Errorf("Config file was not saved at %s", newConfigPath)
	}
}

// Test environment variable overrides
func TestConfigEnvVarOverrides(t *testing.T) {
	// Create a temporary test config file
	testConfig := createTestConfig()
	tempDir, err := ioutil.TempDir("", "beluga-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // clean up after test
	
	configPath := filepath.Join(tempDir, "agents.json")
	if err := writeTestConfig(configPath, testConfig); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	
	// Set environment variables for overrides
	os.Setenv("BELUGA_AGENT_TEST_AGENT_MAX_RETRIES", "10")
	os.Setenv("BELUGA_AGENT_TEST_AGENT_DATA_SOURCE", "override_source")
	defer func() {
		os.Unsetenv("BELUGA_AGENT_TEST_AGENT_MAX_RETRIES")
		os.Unsetenv("BELUGA_AGENT_TEST_AGENT_DATA_SOURCE")
	}()
	
	// Create and initialize config manager
	cm := config.NewConfigManager()
	if err := cm.LoadConfig(configPath); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Check that environment variable overrides were applied
	agentConfig, err := cm.GetAgentConfig("test_agent")
	if err != nil {
		t.Fatalf("Failed to get test_agent config: %v", err)
	}
	
	if agentConfig.MaxRetries != 10 {
		t.Errorf("Environment override for MaxRetries not applied. Expected 10, got %d", agentConfig.MaxRetries)
	}
	
	if dataSource, ok := agentConfig.Settings["data_source"].(string); !ok || dataSource != "override_source" {
		t.Errorf("Environment override for data_source not applied. Expected 'override_source', got %v", agentConfig.Settings["data_source"])
	}
}

// Helper function to create a test config structure
func createTestConfig() *config.AgentModuleConfig {
	return &config.AgentModuleConfig{
		Agents: []*config.AgentConfig{
			{
				Type: "DataFetcherAgent",
				Name: "test_agent",
				Role: "data_fetcher",
				Settings: map[string]interface{}{
					"data_source": "test_source",
					"data_format": "json",
				},
				MaxRetries:   3,
				RetryDelay:   5,
				Dependencies: []string{},
				Description:  "Test agent for unit testing",
			},
			{
				Type: "AnalyzerAgent",
				Name: "analyzer_agent",
				Role: "analyzer",
				Settings: map[string]interface{}{
					"analysis_type": "test",
				},
				MaxRetries:   2,
				RetryDelay:   3,
				Dependencies: []string{"test_agent"},
				Description:  "Test analyzer agent",
			},
		},
		DefaultSettings: map[string]interface{}{
			"max_retries": 3,
			"retry_delay": 5,
			"log_level":   "debug",
		},
		LoggingConfig: map[string]interface{}{
			"console_enabled": true,
			"file_enabled":    false,
			"level":           "debug",
		},
		HealthCheckConfig: map[string]interface{}{
			"enabled":  true,
			"interval": 60,
		},
	}
}

// Helper function to write a test config to a file
func writeTestConfig(path string, config *config.AgentModuleConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}