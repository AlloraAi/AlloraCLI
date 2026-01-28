package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitialize(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "alloracli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test config file
	configFile := filepath.Join(tmpDir, "config.yaml")
	testConfig := `version: "2.0.0"
agents:
  default:
    type: "general"
    model: "gpt-4"
    max_tokens: 4096
    temperature: 0.7
logging:
  level: "info"
  format: "text"
security:
  encryption: true
  audit_logging: true
`

	err = os.WriteFile(configFile, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test initialization
	err = Initialize(configFile, false)
	if err != nil {
		t.Errorf("Initialize() failed: %v", err)
	}

	// Test loading configuration
	cfg, err := Load()
	if err != nil {
		t.Errorf("Load() failed: %v", err)
	}

	// Verify configuration values
	if cfg.Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got '%s'", cfg.Version)
	}

	if len(cfg.Agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(cfg.Agents))
	}

	if cfg.Agents["default"].Type != "general" {
		t.Errorf("Expected agent type 'general', got '%s'", cfg.Agents["default"].Type)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Expected log level 'info', got '%s'", cfg.Logging.Level)
	}

	if !cfg.Security.Encryption {
		t.Error("Expected encryption to be enabled")
	}
}

func TestGetConfigDir(t *testing.T) {
	configDir, err := GetConfigDir()
	if err != nil {
		t.Errorf("GetConfigDir() failed: %v", err)
	}

	if configDir == "" {
		t.Error("GetConfigDir() returned empty string")
	}

	// Check if the path contains expected components
	expectedPath := filepath.Join("alloracli")
	if !filepath.IsAbs(configDir) {
		t.Errorf("GetConfigDir() should return absolute path, got: %s", configDir)
	}

	if !strings.HasSuffix(configDir, expectedPath) {
		t.Errorf("GetConfigDir() should end with 'alloracli', got: %s", configDir)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "alloracli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test configuration
	cfg := &Config{
		Version: "2.0.0",
		Agents: map[string]Agent{
			"test": {
				Type:        "general",
				Model:       "gpt-4",
				MaxTokens:   4096,
				Temperature: 0.7,
			},
		},
		CloudProviders: CloudProviders{
			AWS: AWSConfig{
				Region:  "us-west-2",
				Profile: "default",
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
		Security: SecurityConfig{
			Encryption:   true,
			AuditLogging: true,
		},
	}

	// Test saving configuration
	configFile := filepath.Join(tmpDir, "test-config.yaml")
	err = Save(cfg, configFile)
	if err != nil {
		t.Errorf("Save() failed: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Test loading configuration
	err = Initialize(configFile, false)
	if err != nil {
		t.Errorf("Initialize() failed: %v", err)
	}

	loadedCfg, err := Load()
	if err != nil {
		t.Errorf("Load() failed: %v", err)
	}

	// Verify loaded configuration
	if loadedCfg.Version != cfg.Version {
		t.Errorf("Expected version '%s', got '%s'", cfg.Version, loadedCfg.Version)
	}

	if loadedCfg.Agents["test"].Type != cfg.Agents["test"].Type {
		t.Errorf("Expected agent type '%s', got '%s'", cfg.Agents["test"].Type, loadedCfg.Agents["test"].Type)
	}

	if loadedCfg.CloudProviders.AWS.Region != cfg.CloudProviders.AWS.Region {
		t.Errorf("Expected AWS region '%s', got '%s'", cfg.CloudProviders.AWS.Region, loadedCfg.CloudProviders.AWS.Region)
	}
}

func TestDefaults(t *testing.T) {
	// Test default configuration values
	tmpDir, err := os.MkdirTemp("", "alloracli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize with non-existent config file (should use defaults)
	configFile := filepath.Join(tmpDir, "non-existent.yaml")
	err = Initialize(configFile, false)
	if err != nil {
		t.Errorf("Initialize() with non-existent config failed: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Errorf("Load() failed: %v", err)
	}

	// Check default values
	if cfg.Logging.Level != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", cfg.Logging.Level)
	}

	if cfg.Logging.Format != "text" {
		t.Errorf("Expected default log format 'text', got '%s'", cfg.Logging.Format)
	}

	if !cfg.Security.Encryption {
		t.Error("Expected encryption to be enabled by default")
	}

	if !cfg.Security.AuditLogging {
		t.Error("Expected audit logging to be enabled by default")
	}
}

// BenchmarkLoad benchmarks the configuration loading
func BenchmarkLoad(b *testing.B) {
	// Create a temporary config file
	tmpDir, err := os.MkdirTemp("", "alloracli-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "config.yaml")
	testConfig := `version: "2.0.0"
agents:
  default:
    type: "general"
    model: "gpt-4"
    max_tokens: 4096
    temperature: 0.7
logging:
  level: "info"
  format: "text"
`

	err = os.WriteFile(configFile, []byte(testConfig), 0644)
	if err != nil {
		b.Fatalf("Failed to write test config: %v", err)
	}

	err = Initialize(configFile, false)
	if err != nil {
		b.Fatalf("Initialize() failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Load()
		if err != nil {
			b.Errorf("Load() failed: %v", err)
		}
	}
}
