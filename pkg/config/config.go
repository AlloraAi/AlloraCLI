package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Version        string           `yaml:"version" mapstructure:"version"`
	Agents         map[string]Agent `yaml:"agents" mapstructure:"agents"`
	CloudProviders CloudProviders   `yaml:"cloud_providers" mapstructure:"cloud_providers"`
	Monitoring     MonitoringConfig `yaml:"monitoring" mapstructure:"monitoring"`
	Security       SecurityConfig   `yaml:"security" mapstructure:"security"`
	Plugins        PluginConfig     `yaml:"plugins" mapstructure:"plugins"`
	Logging        LoggingConfig    `yaml:"logging" mapstructure:"logging"`
}

// Agent represents an AI agent configuration
type Agent struct {
	Type        string  `yaml:"type" mapstructure:"type"`
	APIKey      string  `yaml:"api_key" mapstructure:"api_key"`
	Model       string  `yaml:"model" mapstructure:"model"`
	MaxTokens   int     `yaml:"max_tokens" mapstructure:"max_tokens"`
	Temperature float64 `yaml:"temperature" mapstructure:"temperature"`
	Endpoint    string  `yaml:"endpoint,omitempty" mapstructure:"endpoint"`
}

// CloudProviders contains configuration for all cloud providers
type CloudProviders struct {
	AWS   AWSConfig   `yaml:"aws" mapstructure:"aws"`
	Azure AzureConfig `yaml:"azure" mapstructure:"azure"`
	GCP   GCPConfig   `yaml:"gcp" mapstructure:"gcp"`
}

// AWSConfig represents AWS-specific configuration
type AWSConfig struct {
	Region      string `yaml:"region" mapstructure:"region"`
	Profile     string `yaml:"profile" mapstructure:"profile"`
	AccessKeyID string `yaml:"access_key_id,omitempty" mapstructure:"access_key_id"`
	SecretKey   string `yaml:"secret_access_key,omitempty" mapstructure:"secret_access_key"`
}

// AzureConfig represents Azure-specific configuration
type AzureConfig struct {
	SubscriptionID string `yaml:"subscription_id"`
	TenantID       string `yaml:"tenant_id"`
	ClientID       string `yaml:"client_id,omitempty"`
	ClientSecret   string `yaml:"client_secret,omitempty"`
}

// GCPConfig represents GCP-specific configuration
type GCPConfig struct {
	ProjectID          string `yaml:"project_id"`
	Region             string `yaml:"region"`
	ServiceAccountPath string `yaml:"service_account_path,omitempty"`
	ApplicationDefault bool   `yaml:"application_default"`
}

// MonitoringConfig contains monitoring tool configurations
type MonitoringConfig struct {
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Grafana    GrafanaConfig    `yaml:"grafana"`
	DataDog    DataDogConfig    `yaml:"datadog"`
	NewRelic   NewRelicConfig   `yaml:"newrelic"`
}

// PrometheusConfig represents Prometheus configuration
type PrometheusConfig struct {
	Endpoint string `yaml:"endpoint"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// GrafanaConfig represents Grafana configuration
type GrafanaConfig struct {
	Endpoint string `yaml:"endpoint"`
	APIKey   string `yaml:"api_key,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// DataDogConfig represents DataDog configuration
type DataDogConfig struct {
	APIKey string `yaml:"api_key"`
	AppKey string `yaml:"app_key"`
}

// NewRelicConfig represents New Relic configuration
type NewRelicConfig struct {
	APIKey    string `yaml:"api_key"`
	AccountID string `yaml:"account_id"`
}

// SecurityConfig contains security-related settings
type SecurityConfig struct {
	Encryption     bool   `yaml:"encryption" mapstructure:"encryption"`
	AuditLogging   bool   `yaml:"audit_logging" mapstructure:"audit_logging"`
	KeyManagement  string `yaml:"key_management" mapstructure:"key_management"`
	ComplianceMode string `yaml:"compliance_mode" mapstructure:"compliance_mode"`
}

// PluginConfig contains plugin-related settings
type PluginConfig struct {
	Directory      string   `yaml:"directory"`
	AutoUpdate     bool     `yaml:"auto_update"`
	AllowedSources []string `yaml:"allowed_sources"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level    string `yaml:"level" mapstructure:"level"`
	Format   string `yaml:"format" mapstructure:"format"`
	Output   string `yaml:"output" mapstructure:"output"`
	Rotate   bool   `yaml:"rotate" mapstructure:"rotate"`
	MaxSize  int    `yaml:"max_size" mapstructure:"max_size"`
	MaxAge   int    `yaml:"max_age" mapstructure:"max_age"`
	MaxFiles int    `yaml:"max_files" mapstructure:"max_files"`
}

// Initialize initializes the configuration system
func Initialize(configFile string, verbose bool) error {
	// Set config file path
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Set default config paths
		configDir, err := GetConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %w", err)
		}

		viper.AddConfigPath(configDir)
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Environment variables
	viper.SetEnvPrefix("ALLORA")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// If a specific config file was provided and it doesn't exist,
			// treat it as a missing config file error instead of a read error
			if configFile != "" {
				if os.IsNotExist(err) {
					// Continue without error, just use defaults
				} else {
					return fmt.Errorf("failed to read config file: %w", err)
				}
			} else {
				return fmt.Errorf("failed to read config file: %w", err)
			}
		}
	}

	return nil
}

// Load loads the configuration from file
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

// Save saves the configuration to file
func Save(cfg *Config, configFile string) error {
	if configFile == "" {
		configDir, err := GetConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %w", err)
		}
		configFile = filepath.Join(configDir, "config.yaml")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Display displays the configuration in the specified format
func Display(cfg *Config, format string) error {
	switch format {
	case "json":
		return displayJSON(cfg)
	case "yaml":
		return displayYAML(cfg)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			// Fallback to USERPROFILE if APPDATA is not set
			userProfile := os.Getenv("USERPROFILE")
			if userProfile != "" {
				configDir = filepath.Join(userProfile, "AppData", "Roaming", "alloracli")
			} else {
				return "", fmt.Errorf("neither APPDATA nor USERPROFILE environment variables are set")
			}
		} else {
			configDir = filepath.Join(configDir, "alloracli")
		}
	default:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config", "alloracli")
	}

	return configDir, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "text")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.rotate", true)
	viper.SetDefault("logging.max_size", 100)
	viper.SetDefault("logging.max_age", 30)
	viper.SetDefault("logging.max_files", 10)

	// Security defaults
	viper.SetDefault("security.encryption", true)
	viper.SetDefault("security.audit_logging", true)
	viper.SetDefault("security.key_management", "local")

	// Plugin defaults
	viper.SetDefault("plugins.auto_update", false)
	viper.SetDefault("plugins.allowed_sources", []string{"github.com", "registry.alloraai.com"})

	// Cloud provider defaults
	viper.SetDefault("cloud_providers.aws.region", "us-west-2")
	viper.SetDefault("cloud_providers.aws.profile", "default")
	viper.SetDefault("cloud_providers.gcp.region", "us-central1")
	viper.SetDefault("cloud_providers.gcp.application_default", true)

	// Monitoring defaults
	viper.SetDefault("monitoring.prometheus.endpoint", "http://localhost:9090")
	viper.SetDefault("monitoring.grafana.endpoint", "http://localhost:3000")
}

// displayJSON displays configuration in JSON format
func displayJSON(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// displayYAML displays configuration in YAML format
func displayYAML(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
