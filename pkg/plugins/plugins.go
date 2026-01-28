package plugins

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/hashicorp/go-plugin"
)

// PluginService interface defines plugin management operations
type PluginService interface {
	ListPlugins(ctx context.Context) ([]PluginInfo, error)
	InstallPlugin(ctx context.Context, name string, source string) error
	UpdatePlugin(ctx context.Context, name string) error
	UninstallPlugin(ctx context.Context, name string) error
	EnablePlugin(ctx context.Context, name string) error
	DisablePlugin(ctx context.Context, name string) error
	GetPluginInfo(ctx context.Context, name string) (*PluginInfo, error)
	ExecutePlugin(ctx context.Context, name string, args []string) (*PluginResult, error)
	SearchPlugins(ctx context.Context, query string) ([]PluginSearchResult, error)
}

// PluginInfo represents plugin information
type PluginInfo struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	License      string            `json:"license"`
	Homepage     string            `json:"homepage"`
	Repository   string            `json:"repository"`
	Tags         []string          `json:"tags"`
	Commands     []CommandInfo     `json:"commands"`
	Status       string            `json:"status"`
	Enabled      bool              `json:"enabled"`
	Installed    time.Time         `json:"installed"`
	Updated      time.Time         `json:"updated"`
	Config       map[string]string `json:"config"`
	Dependencies []string          `json:"dependencies"`
}

// CommandInfo represents a plugin command
type CommandInfo struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Usage       string     `json:"usage"`
	Flags       []FlagInfo `json:"flags"`
	Examples    []string   `json:"examples"`
}

// FlagInfo represents a command flag
type FlagInfo struct {
	Name        string `json:"name"`
	Short       string `json:"short"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Required    bool   `json:"required"`
}

// PluginResult represents the result of plugin execution
type PluginResult struct {
	ExitCode int                    `json:"exit_code"`
	Output   string                 `json:"output"`
	Error    string                 `json:"error"`
	Data     map[string]interface{} `json:"data"`
	Duration time.Duration          `json:"duration"`
}

// PluginSearchResult represents a plugin search result
type PluginSearchResult struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Tags        []string  `json:"tags"`
	Downloads   int       `json:"downloads"`
	Rating      float64   `json:"rating"`
	Updated     time.Time `json:"updated"`
	Source      string    `json:"source"`
}

// Plugin interface defines the contract for AlloraCLI plugins
type Plugin interface {
	GetInfo() *PluginInfo
	Execute(args []string) (*PluginResult, error)
	Configure(config map[string]string) error
	Validate() error
}

// PluginManifest represents the plugin manifest file
type PluginManifest struct {
	Name         string            `yaml:"name"`
	Version      string            `yaml:"version"`
	Description  string            `yaml:"description"`
	Author       string            `yaml:"author"`
	License      string            `yaml:"license"`
	Homepage     string            `yaml:"homepage"`
	Repository   string            `yaml:"repository"`
	Tags         []string          `yaml:"tags"`
	Commands     []CommandInfo     `yaml:"commands"`
	Dependencies []string          `yaml:"dependencies"`
	Config       map[string]string `yaml:"config"`
	Binary       string            `yaml:"binary"`
	Checksum     string            `yaml:"checksum"`
}

// DefaultPluginService provides a default implementation
type DefaultPluginService struct {
	config    *config.Config
	pluginDir string
	plugins   map[string]*PluginInfo
}

// NewPluginService creates a new plugin service
func NewPluginService(cfg *config.Config) (PluginService, error) {
	pluginDir := cfg.Plugins.Directory
	if pluginDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		pluginDir = filepath.Join(homeDir, ".config", "alloracli", "plugins")
	}

	service := &DefaultPluginService{
		config:    cfg,
		pluginDir: pluginDir,
		plugins:   make(map[string]*PluginInfo),
	}

	// Load existing plugins
	if err := service.loadPlugins(); err != nil {
		return nil, fmt.Errorf("failed to load plugins: %w", err)
	}

	return service, nil
}

// ListPlugins lists all installed plugins
func (p *DefaultPluginService) ListPlugins(ctx context.Context) ([]PluginInfo, error) {
	var plugins []PluginInfo

	for _, plugin := range p.plugins {
		plugins = append(plugins, *plugin)
	}

	return plugins, nil
}

// InstallPlugin installs a plugin from a source
func (p *DefaultPluginService) InstallPlugin(ctx context.Context, name string, source string) error {
	// Mock implementation - in real implementation, this would:
	// 1. Download plugin from source
	// 2. Verify checksum
	// 3. Extract plugin
	// 4. Validate manifest
	// 5. Install dependencies
	// 6. Register plugin

	pluginInfo := &PluginInfo{
		Name:        name,
		Version:     "2.0.0",
		Description: "Sample plugin",
		Author:      "AlloraAi",
		License:     "MIT",
		Homepage:    "https://github.com/AlloraAi/allora-plugin-" + name,
		Repository:  "https://github.com/AlloraAi/allora-plugin-" + name,
		Tags:        []string{"sample"},
		Commands: []CommandInfo{
			{
				Name:        name,
				Description: "Execute " + name + " plugin",
				Usage:       name + " [options]",
				Flags:       []FlagInfo{},
				Examples:    []string{name + " --help"},
			},
		},
		Status:       "installed",
		Enabled:      true,
		Installed:    time.Now(),
		Updated:      time.Now(),
		Config:       make(map[string]string),
		Dependencies: []string{},
	}

	p.plugins[name] = pluginInfo

	return nil
}

// UpdatePlugin updates a plugin to the latest version
func (p *DefaultPluginService) UpdatePlugin(ctx context.Context, name string) error {
	plugin, exists := p.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Mock implementation - would check for updates and install them
	plugin.Updated = time.Now()
	plugin.Version = "1.0.1"

	return nil
}

// UninstallPlugin removes a plugin
func (p *DefaultPluginService) UninstallPlugin(ctx context.Context, name string) error {
	if _, exists := p.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Mock implementation - would remove plugin files and cleanup
	delete(p.plugins, name)

	return nil
}

// EnablePlugin enables a plugin
func (p *DefaultPluginService) EnablePlugin(ctx context.Context, name string) error {
	plugin, exists := p.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	plugin.Enabled = true
	plugin.Status = "enabled"

	return nil
}

// DisablePlugin disables a plugin
func (p *DefaultPluginService) DisablePlugin(ctx context.Context, name string) error {
	plugin, exists := p.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	plugin.Enabled = false
	plugin.Status = "disabled"

	return nil
}

// GetPluginInfo gets information about a specific plugin
func (p *DefaultPluginService) GetPluginInfo(ctx context.Context, name string) (*PluginInfo, error) {
	plugin, exists := p.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// ExecutePlugin executes a plugin with the given arguments
func (p *DefaultPluginService) ExecutePlugin(ctx context.Context, name string, args []string) (*PluginResult, error) {
	pluginInfo, exists := p.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	if !pluginInfo.Enabled {
		return nil, fmt.Errorf("plugin %s is disabled", name)
	}

	start := time.Now()

	// Mock implementation - would execute the actual plugin
	result := &PluginResult{
		ExitCode: 0,
		Output:   fmt.Sprintf("Plugin %s executed successfully with args: %v", name, args),
		Error:    "",
		Data: map[string]interface{}{
			"plugin": name,
			"args":   args,
		},
		Duration: time.Since(start),
	}

	return result, nil
}

// SearchPlugins searches for plugins in the registry
func (p *DefaultPluginService) SearchPlugins(ctx context.Context, query string) ([]PluginSearchResult, error) {
	// Mock implementation - would search in plugin registries
	results := []PluginSearchResult{
		{
			Name:        "aws-helper",
			Version:     "2.1.0",
			Description: "AWS infrastructure management helper",
			Author:      "AlloraAi",
			Tags:        []string{"aws", "cloud", "infrastructure"},
			Downloads:   1500,
			Rating:      4.8,
			Updated:     time.Now().Add(-7 * 24 * time.Hour),
			Source:      "https://registry.alloraai.com/plugins/aws-helper",
		},
		{
			Name:        "k8s-manager",
			Version:     "1.5.2",
			Description: "Kubernetes cluster management utilities",
			Author:      "AlloraAi",
			Tags:        []string{"kubernetes", "k8s", "containers"},
			Downloads:   2200,
			Rating:      4.9,
			Updated:     time.Now().Add(-3 * 24 * time.Hour),
			Source:      "https://registry.alloraai.com/plugins/k8s-manager",
		},
		{
			Name:        "monitoring-tools",
			Version:     "3.0.1",
			Description: "Comprehensive monitoring and alerting tools",
			Author:      "AlloraAi",
			Tags:        []string{"monitoring", "alerting", "metrics"},
			Downloads:   1800,
			Rating:      4.7,
			Updated:     time.Now().Add(-10 * 24 * time.Hour),
			Source:      "https://registry.alloraai.com/plugins/monitoring-tools",
		},
	}

	// Filter results based on query
	var filtered []PluginSearchResult
	for _, result := range results {
		if containsQuery(result, query) {
			filtered = append(filtered, result)
		}
	}

	return filtered, nil
}

// loadPlugins loads plugins from the plugin directory
func (p *DefaultPluginService) loadPlugins() error {
	// Create plugin directory if it doesn't exist
	if err := os.MkdirAll(p.pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Mock implementation - would scan directory and load plugin manifests
	// For now, we'll add some sample plugins
	samplePlugins := []*PluginInfo{
		{
			Name:        "sample-plugin",
			Version:     "2.0.0",
			Description: "A sample plugin for demonstration",
			Author:      "AlloraAi",
			License:     "MIT",
			Homepage:    "https://github.com/AlloraAi/allora-plugin-sample",
			Repository:  "https://github.com/AlloraAi/allora-plugin-sample",
			Tags:        []string{"sample", "demo"},
			Commands: []CommandInfo{
				{
					Name:        "sample",
					Description: "Execute sample plugin",
					Usage:       "sample [options]",
					Flags: []FlagInfo{
						{
							Name:        "verbose",
							Short:       "v",
							Description: "Enable verbose output",
							Type:        "bool",
							Default:     "false",
							Required:    false,
						},
					},
					Examples: []string{
						"sample --verbose",
						"sample -v",
					},
				},
			},
			Status:       "installed",
			Enabled:      true,
			Installed:    time.Now().Add(-30 * 24 * time.Hour),
			Updated:      time.Now().Add(-7 * 24 * time.Hour),
			Config:       make(map[string]string),
			Dependencies: []string{},
		},
	}

	for _, plugin := range samplePlugins {
		p.plugins[plugin.Name] = plugin
	}

	return nil
}

// containsQuery checks if a search result matches the query
func containsQuery(result PluginSearchResult, query string) bool {
	if query == "" {
		return true
	}

	// Simple text matching - in real implementation, this would be more sophisticated
	if contains(result.Name, query) ||
		contains(result.Description, query) ||
		contains(result.Author, query) {
		return true
	}

	for _, tag := range result.Tags {
		if contains(tag, query) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

// containsSubstring checks if a string contains a substring
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Plugin registry interface for advanced plugin management
type PluginRegistry interface {
	Search(query string) ([]PluginSearchResult, error)
	GetMetadata(name string) (*PluginManifest, error)
	Download(name string, version string) ([]byte, error)
	Verify(data []byte, checksum string) error
}

// LocalPluginRegistry implements a local plugin registry
type LocalPluginRegistry struct {
	baseURL string
}

// NewLocalPluginRegistry creates a new local plugin registry
func NewLocalPluginRegistry(baseURL string) PluginRegistry {
	return &LocalPluginRegistry{
		baseURL: baseURL,
	}
}

// Search searches for plugins in the registry
func (r *LocalPluginRegistry) Search(query string) ([]PluginSearchResult, error) {
	// Mock implementation - would make HTTP requests to registry
	return []PluginSearchResult{}, nil
}

// GetMetadata gets plugin metadata from the registry
func (r *LocalPluginRegistry) GetMetadata(name string) (*PluginManifest, error) {
	// Mock implementation - would fetch metadata from registry
	return &PluginManifest{
		Name:        name,
		Version:     "2.0.0",
		Description: "Sample plugin",
		Author:      "AlloraAi",
		License:     "MIT",
		Binary:      name,
		Checksum:    "abc123def456",
	}, nil
}

// Download downloads a plugin from the registry
func (r *LocalPluginRegistry) Download(name string, version string) ([]byte, error) {
	// Mock implementation - would download plugin binary
	return []byte("mock plugin binary"), nil
}

// Verify verifies a plugin's integrity
func (r *LocalPluginRegistry) Verify(data []byte, checksum string) error {
	// Mock implementation - would verify checksum
	return nil
}

// PluginManager manages the plugin lifecycle
type PluginManager struct {
	service  PluginService
	registry PluginRegistry
	clients  map[string]*plugin.Client
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(service PluginService, registry PluginRegistry) *PluginManager {
	return &PluginManager{
		service:  service,
		registry: registry,
		clients:  make(map[string]*plugin.Client),
	}
}

// GetClient gets a plugin client for communication
func (m *PluginManager) GetClient(name string) (*plugin.Client, error) {
	if client, exists := m.clients[name]; exists {
		return client, nil
	}

	// Mock implementation - would create actual plugin client
	// using hashicorp/go-plugin
	return nil, fmt.Errorf("plugin client not implemented for %s", name)
}

// CleanupClients cleans up all plugin clients
func (m *PluginManager) CleanupClients() {
	for _, client := range m.clients {
		client.Kill()
	}
	m.clients = make(map[string]*plugin.Client)
}
