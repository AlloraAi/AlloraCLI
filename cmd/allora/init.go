package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize AlloraCLI configuration and authentication",
		Long:  `Initialize AlloraCLI by setting up configuration files, authentication, and basic agent setup.`,
		RunE:  runInit,
	}

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	utils.PrintBanner()

	fmt.Println("🚀 Welcome to AlloraCLI!")
	fmt.Println("Let's set up your AI-powered IT infrastructure management CLI.")

	// Check if already initialized
	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configFile); err == nil {
		prompt := promptui.Prompt{
			Label:     "AlloraCLI is already initialized. Do you want to reinitialize",
			IsConfirm: true,
		}

		if _, err := prompt.Run(); err != nil {
			fmt.Println("Initialization cancelled.")
			return nil
		}
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Interactive setup
	cfg := &config.Config{
		Version: "2.0.0",
		Agents:  make(map[string]config.Agent),
		CloudProviders: config.CloudProviders{
			AWS:   config.AWSConfig{},
			Azure: config.AzureConfig{},
			GCP:   config.GCPConfig{},
		},
		Monitoring: config.MonitoringConfig{},
		Security: config.SecurityConfig{
			Encryption:   true,
			AuditLogging: true,
		},
	}

	// Setup default agent
	if err := setupDefaultAgent(cfg); err != nil {
		return fmt.Errorf("failed to setup default agent: %w", err)
	}

	// Setup cloud providers (optional)
	if err := setupCloudProviders(cfg); err != nil {
		return fmt.Errorf("failed to setup cloud providers: %w", err)
	}

	// Setup monitoring (optional)
	if err := setupMonitoring(cfg); err != nil {
		return fmt.Errorf("failed to setup monitoring: %w", err)
	}

	// Save configuration
	if err := config.Save(cfg, configFile); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Success message
	color.Green("✅ AlloraCLI has been successfully initialized!")
	fmt.Printf("\nConfiguration saved to: %s\n", configFile)
	fmt.Println("\nYou can now start using AlloraCLI:")
	fmt.Println("  allora ask \"What's the status of my infrastructure?\"")
	fmt.Println("  allora monitor status")
	fmt.Println("  allora troubleshoot --help")
	fmt.Println("\nFor more information, visit: https://docs.alloraai.com")

	return nil
}

func setupDefaultAgent(cfg *config.Config) error {
	fmt.Println("\n🤖 Setting up your first AI agent...")

	// Agent name
	prompt := promptui.Prompt{
		Label:   "Agent name",
		Default: "infra-assistant",
	}
	name, err := prompt.Run()
	if err != nil {
		return err
	}

	// Agent type
	selectPrompt := promptui.Select{
		Label: "Agent type",
		Items: []string{"general", "aws", "azure", "gcp", "kubernetes", "monitoring"},
	}
	_, agentType, err := selectPrompt.Run()
	if err != nil {
		return err
	}

	// API Key (optional for now)
	apiKeyPrompt := promptui.Prompt{
		Label: "API Key (optional, can be set later via environment variable ALLORA_API_KEY)",
		Mask:  '*',
	}
	apiKey, err := apiKeyPrompt.Run()
	if err != nil && err != promptui.ErrInterrupt {
		return err
	}

	// Model selection
	modelSelect := promptui.Select{
		Label: "AI Model",
		Items: []string{"gpt-4", "gpt-3.5-turbo", "claude-3", "gemini-pro"},
	}
	_, model, err := modelSelect.Run()
	if err != nil {
		return err
	}

	cfg.Agents[name] = config.Agent{
		Type:        agentType,
		APIKey:      apiKey,
		Model:       model,
		MaxTokens:   2048,
		Temperature: 0.7,
	}

	return nil
}

func setupCloudProviders(cfg *config.Config) error {
	fmt.Println("\n☁️  Cloud Provider Setup (optional)")

	prompt := promptui.Prompt{
		Label:     "Do you want to configure cloud providers now",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		return nil // Skip cloud provider setup
	}

	// AWS Setup
	awsPrompt := promptui.Prompt{
		Label:     "Configure AWS",
		IsConfirm: true,
	}
	if _, err := awsPrompt.Run(); err == nil {
		regionPrompt := promptui.Prompt{
			Label:   "AWS Region",
			Default: "us-west-2",
		}
		region, err := regionPrompt.Run()
		if err != nil {
			return err
		}

		profilePrompt := promptui.Prompt{
			Label:   "AWS Profile",
			Default: "default",
		}
		profile, err := profilePrompt.Run()
		if err != nil {
			return err
		}

		cfg.CloudProviders.AWS = config.AWSConfig{
			Region:  region,
			Profile: profile,
		}
	}

	// Azure Setup
	azurePrompt := promptui.Prompt{
		Label:     "Configure Azure",
		IsConfirm: true,
	}
	if _, err := azurePrompt.Run(); err == nil {
		subscriptionPrompt := promptui.Prompt{
			Label: "Azure Subscription ID",
		}
		subscriptionID, err := subscriptionPrompt.Run()
		if err != nil {
			return err
		}

		tenantPrompt := promptui.Prompt{
			Label: "Azure Tenant ID",
		}
		tenantID, err := tenantPrompt.Run()
		if err != nil {
			return err
		}

		cfg.CloudProviders.Azure = config.AzureConfig{
			SubscriptionID: subscriptionID,
			TenantID:       tenantID,
		}
	}

	return nil
}

func setupMonitoring(cfg *config.Config) error {
	fmt.Println("\n📊 Monitoring Setup (optional)")

	prompt := promptui.Prompt{
		Label:     "Do you want to configure monitoring tools now",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		return nil // Skip monitoring setup
	}

	// Prometheus
	prometheusPrompt := promptui.Prompt{
		Label:     "Configure Prometheus",
		IsConfirm: true,
	}
	if _, err := prometheusPrompt.Run(); err == nil {
		endpointPrompt := promptui.Prompt{
			Label:   "Prometheus Endpoint",
			Default: "http://localhost:9090",
		}
		endpoint, err := endpointPrompt.Run()
		if err != nil {
			return err
		}

		cfg.Monitoring.Prometheus = config.PrometheusConfig{
			Endpoint: endpoint,
		}
	}

	// Grafana
	grafanaPrompt := promptui.Prompt{
		Label:     "Configure Grafana",
		IsConfirm: true,
	}
	if _, err := grafanaPrompt.Run(); err == nil {
		endpointPrompt := promptui.Prompt{
			Label:   "Grafana Endpoint",
			Default: "http://localhost:3000",
		}
		endpoint, err := endpointPrompt.Run()
		if err != nil {
			return err
		}

		cfg.Monitoring.Grafana = config.GrafanaConfig{
			Endpoint: endpoint,
		}
	}

	return nil
}
