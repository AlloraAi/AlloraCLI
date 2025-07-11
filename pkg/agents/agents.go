package agents

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/go-resty/resty/v2"
)

// Agent represents an AI agent interface
type Agent interface {
	GetName() string
	GetType() string
	Query(ctx context.Context, query *Query) (*Response, error)
	GetCapabilities() []string
	GetStatus() *AgentStatus
	GetConfiguration() *AgentConfig
	UpdateConfiguration(config *AgentConfig) error
	Start() error
	Stop() error
	IsHealthy() bool
}

// Query represents a query to an AI agent
type Query struct {
	Text    string                 `json:"text"`
	Context map[string]interface{} `json:"context"`
}

// Response represents an AI agent response
type Response struct {
	Text       string                 `json:"text"`
	Content    string                 `json:"content"`
	Type       string                 `json:"type"`
	Confidence float64                `json:"confidence"`
	Metadata   map[string]interface{} `json:"metadata"`
	Suggestions []string              `json:"suggestions"`
	Actions    []Action               `json:"actions"`
	Timestamp  time.Time              `json:"timestamp"`
}

// Action represents an actionable item from the AI response
type Action struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Command     string            `json:"command"`
	Parameters  map[string]interface{} `json:"parameters"`
	Risk        string            `json:"risk"`
}

// AgentStatus represents the status of an AI agent
type AgentStatus struct {
	State        string    `json:"state"`
	LastActivity time.Time `json:"last_activity"`
	Health       string    `json:"health"`
	Uptime       time.Duration `json:"uptime"`
}

// AgentConfig represents the configuration of an AI agent
type AgentConfig struct {
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	APIKey      string  `json:"api_key"`
	Endpoint    string  `json:"endpoint"`
}

// AgentManager manages multiple AI agents
type AgentManager struct {
	agents map[string]Agent
	mutex  sync.RWMutex
}

// NewAgentManager creates a new agent manager
func NewAgentManager() *AgentManager {
	return &AgentManager{
		agents: make(map[string]Agent),
	}
}

// AddAgent adds an agent to the manager
func (m *AgentManager) AddAgent(agent Agent) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.agents[agent.GetName()] = agent
	return nil
}

// GetAgent retrieves an agent by name
func (m *AgentManager) GetAgent(name string) (Agent, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	agent, exists := m.agents[name]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", name)
	}
	return agent, nil
}

// RemoveAgent removes an agent from the manager
func (m *AgentManager) RemoveAgent(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.agents[name]; !exists {
		return fmt.Errorf("agent not found: %s", name)
	}
	
	delete(m.agents, name)
	return nil
}

// ListAgents returns a list of all agents
func (m *AgentManager) ListAgents() []Agent {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	agents := make([]Agent, 0, len(m.agents))
	for _, agent := range m.agents {
		agents = append(agents, agent)
	}
	return agents
}

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	name     string
	config   config.Agent
	client   *resty.Client
	context  context.Context
	status   *AgentStatus
	agentConfig *AgentConfig
}

// GetName returns the agent name
func (b *BaseAgent) GetName() string {
	return b.name
}

// GetType returns the agent type
func (b *BaseAgent) GetType() string {
	return b.config.Type
}

// Query processes a query and returns a response
func (b *BaseAgent) Query(ctx context.Context, query *Query) (*Response, error) {
	// This is a base implementation - specific agents should override this
	return &Response{
		Text:       fmt.Sprintf("Processed query: %s", query.Text),
		Confidence: 0.8,
		Metadata: map[string]interface{}{
			"agent_type": b.GetType(),
			"timestamp":  time.Now().UTC(),
		},
		Timestamp: time.Now().UTC(),
	}, nil
}

// GetCapabilities returns the agent capabilities
func (b *BaseAgent) GetCapabilities() []string {
	return []string{"general", "chat", "help"}
}

// GetStatus returns the agent status
func (b *BaseAgent) GetStatus() *AgentStatus {
	if b.status == nil {
		b.status = &AgentStatus{
			State:        "idle",
			LastActivity: time.Now().UTC(),
			Health:       "healthy",
		}
	}
	return b.status
}

// GetConfiguration returns the agent configuration
func (b *BaseAgent) GetConfiguration() *AgentConfig {
	if b.agentConfig == nil {
		b.agentConfig = &AgentConfig{
			Model:       b.config.Model,
			MaxTokens:   b.config.MaxTokens,
			Temperature: b.config.Temperature,
			APIKey:      b.config.APIKey,
			Endpoint:    b.config.Endpoint,
		}
	}
	return b.agentConfig
}

// UpdateConfiguration updates the agent configuration
func (b *BaseAgent) UpdateConfiguration(config *AgentConfig) error {
	b.agentConfig = config
	b.config.Model = config.Model
	b.config.MaxTokens = config.MaxTokens
	b.config.Temperature = config.Temperature
	b.config.APIKey = config.APIKey
	b.config.Endpoint = config.Endpoint
	return nil
}

// Start starts the agent
func (b *BaseAgent) Start() error {
	if b.status == nil {
		b.status = &AgentStatus{}
	}
	b.status.State = "running"
	b.status.LastActivity = time.Now().UTC()
	return nil
}

// Stop stops the agent
func (b *BaseAgent) Stop() error {
	if b.status == nil {
		b.status = &AgentStatus{}
	}
	b.status.State = "stopped"
	return nil
}

// IsHealthy returns whether the agent is healthy
func (b *BaseAgent) IsHealthy() bool {
	return b.GetStatus().Health == "healthy"
}

// GeneralAgent handles general IT infrastructure questions
type GeneralAgent struct {
	*BaseAgent
}

// AWSAgent specializes in AWS-related queries
type AWSAgent struct {
	*BaseAgent
}

// AzureAgent specializes in Azure-related queries
type AzureAgent struct {
	*BaseAgent
}

// GCPAgent specializes in GCP-related queries
type GCPAgent struct {
	*BaseAgent
}

// KubernetesAgent specializes in Kubernetes-related queries
type KubernetesAgent struct {
	*BaseAgent
}

// MonitoringAgent specializes in monitoring and observability
type MonitoringAgent struct {
	*BaseAgent
}

// NewAgent creates a new agent based on the configuration
func NewAgent(cfg config.Agent) (Agent, error) {
	// Check if this should be an OpenAI agent
	if cfg.APIKey != "" && (cfg.Model == "gpt-4" || cfg.Model == "gpt-3.5-turbo" || cfg.Model == "gpt-4-turbo") {
		return NewOpenAIAgent(cfg, cfg.Type)
	}

	// Fallback to mock agents for development/testing
	baseAgent := &BaseAgent{
		name:    cfg.Type,
		config:  cfg,
		client:  resty.New(),
		context: context.Background(),
	}

	// Configure HTTP client
	baseAgent.client.SetTimeout(30 * time.Second)
	baseAgent.client.SetRetryCount(3)

	// Set API key if provided
	if cfg.APIKey != "" {
		baseAgent.client.SetHeader("Authorization", "Bearer "+cfg.APIKey)
	}

	// Create specific agent based on type
	switch cfg.Type {
	case "general":
		return &GeneralAgent{BaseAgent: baseAgent}, nil
	case "aws":
		return &AWSAgent{BaseAgent: baseAgent}, nil
	case "azure":
		return &AzureAgent{BaseAgent: baseAgent}, nil
	case "gcp":
		return &GCPAgent{BaseAgent: baseAgent}, nil
	case "kubernetes":
		return &KubernetesAgent{BaseAgent: baseAgent}, nil
	case "monitoring":
		return &MonitoringAgent{BaseAgent: baseAgent}, nil
	default:
		return &GeneralAgent{BaseAgent: baseAgent}, nil
	}
}

// ProcessQuery processes a query for the general agent
func (g *GeneralAgent) ProcessQuery(query string) (*Response, error) {
	// This is a mock implementation
	// In a real implementation, you would call the actual AI service
	response := &Response{
		Content:    g.generateMockResponse(query),
		Type:       "text",
		Confidence: 0.85,
		Metadata: map[string]interface{}{
			"agent_type": "general",
			"model":      g.config.Model,
		},
		Suggestions: g.generateSuggestions(query),
		Actions:     g.generateActions(query),
		Timestamp:   time.Now(),
	}

	return response, nil
}

// GetCapabilities returns the capabilities of the general agent
func (g *GeneralAgent) GetCapabilities() []string {
	return []string{
		"Infrastructure monitoring",
		"System troubleshooting",
		"Performance analysis",
		"Security recommendations",
		"Best practices guidance",
	}
}

// GetType returns the agent type
func (g *GeneralAgent) GetType() string {
	return "general"
}

// GetModel returns the AI model being used
func (g *GeneralAgent) GetModel() string {
	return g.config.Model
}

// ProcessQuery processes a query for the AWS agent
func (a *AWSAgent) ProcessQuery(query string) (*Response, error) {
	response := &Response{
		Content:    a.generateAWSResponse(query),
		Type:       "text",
		Confidence: 0.90,
		Metadata: map[string]interface{}{
			"agent_type": "aws",
			"model":      a.config.Model,
		},
		Suggestions: a.generateAWSSuggestions(query),
		Actions:     a.generateAWSActions(query),
		Timestamp:   time.Now(),
	}

	return response, nil
}

// GetCapabilities returns the capabilities of the AWS agent
func (a *AWSAgent) GetCapabilities() []string {
	return []string{
		"EC2 instance management",
		"S3 bucket operations",
		"RDS database monitoring",
		"Lambda function analysis",
		"Cost optimization",
		"Security best practices",
	}
}

// GetType returns the agent type
func (a *AWSAgent) GetType() string {
	return "aws"
}

// GetModel returns the AI model being used
func (a *AWSAgent) GetModel() string {
	return a.config.Model
}

// ProcessQuery processes a query for the Azure agent
func (az *AzureAgent) ProcessQuery(query string) (*Response, error) {
	response := &Response{
		Content:    az.generateAzureResponse(query),
		Type:       "text",
		Confidence: 0.88,
		Metadata: map[string]interface{}{
			"agent_type": "azure",
			"model":      az.config.Model,
		},
		Suggestions: az.generateAzureSuggestions(query),
		Actions:     az.generateAzureActions(query),
		Timestamp:   time.Now(),
	}

	return response, nil
}

// GetCapabilities returns the capabilities of the Azure agent
func (az *AzureAgent) GetCapabilities() []string {
	return []string{
		"Virtual machine management",
		"Storage account operations",
		"Azure SQL monitoring",
		"Function app analysis",
		"Cost management",
		"Security recommendations",
	}
}

// GetType returns the agent type
func (az *AzureAgent) GetType() string {
	return "azure"
}

// GetModel returns the AI model being used
func (az *AzureAgent) GetModel() string {
	return az.config.Model
}

// ProcessQuery processes a query for the GCP agent
func (g *GCPAgent) ProcessQuery(query string) (*Response, error) {
	response := &Response{
		Content:    g.generateGCPResponse(query),
		Type:       "text",
		Confidence: 0.87,
		Metadata: map[string]interface{}{
			"agent_type": "gcp",
			"model":      g.config.Model,
		},
		Suggestions: g.generateGCPSuggestions(query),
		Actions:     g.generateGCPActions(query),
		Timestamp:   time.Now(),
	}

	return response, nil
}

// GetCapabilities returns the capabilities of the GCP agent
func (g *GCPAgent) GetCapabilities() []string {
	return []string{
		"Compute Engine management",
		"Cloud Storage operations",
		"Cloud SQL monitoring",
		"Cloud Functions analysis",
		"Billing optimization",
		"Security recommendations",
	}
}

// GetType returns the agent type
func (g *GCPAgent) GetType() string {
	return "gcp"
}

// GetModel returns the AI model being used
func (g *GCPAgent) GetModel() string {
	return g.config.Model
}

// ProcessQuery processes a query for the Kubernetes agent
func (k *KubernetesAgent) ProcessQuery(query string) (*Response, error) {
	response := &Response{
		Content:    k.generateK8sResponse(query),
		Type:       "text",
		Confidence: 0.89,
		Metadata: map[string]interface{}{
			"agent_type": "kubernetes",
			"model":      k.config.Model,
		},
		Suggestions: k.generateK8sSuggestions(query),
		Actions:     k.generateK8sActions(query),
		Timestamp:   time.Now(),
	}

	return response, nil
}

// GetCapabilities returns the capabilities of the Kubernetes agent
func (k *KubernetesAgent) GetCapabilities() []string {
	return []string{
		"Pod management",
		"Service discovery",
		"Deployment strategies",
		"Resource optimization",
		"Network policies",
		"Security scanning",
	}
}

// GetType returns the agent type
func (k *KubernetesAgent) GetType() string {
	return "kubernetes"
}

// GetModel returns the AI model being used
func (k *KubernetesAgent) GetModel() string {
	return k.config.Model
}

// ProcessQuery processes a query for the monitoring agent
func (m *MonitoringAgent) ProcessQuery(query string) (*Response, error) {
	response := &Response{
		Content:    m.generateMonitoringResponse(query),
		Type:       "text",
		Confidence: 0.91,
		Metadata: map[string]interface{}{
			"agent_type": "monitoring",
			"model":      m.config.Model,
		},
		Suggestions: m.generateMonitoringSuggestions(query),
		Actions:     m.generateMonitoringActions(query),
		Timestamp:   time.Now(),
	}

	return response, nil
}

// GetCapabilities returns the capabilities of the monitoring agent
func (m *MonitoringAgent) GetCapabilities() []string {
	return []string{
		"Metrics analysis",
		"Alert management",
		"Dashboard creation",
		"SLO monitoring",
		"Anomaly detection",
		"Performance tuning",
	}
}

// GetType returns the agent type
func (m *MonitoringAgent) GetType() string {
	return "monitoring"
}

// GetModel returns the AI model being used
func (m *MonitoringAgent) GetModel() string {
	return m.config.Model
}

// Mock response generators (in a real implementation, these would call actual AI services)
func (g *GeneralAgent) generateMockResponse(query string) string {
	return fmt.Sprintf("Based on your query about '%s', here's my analysis:\n\n"+
		"This appears to be a general IT infrastructure question. I can help you with:\n"+
		"- System monitoring and health checks\n"+
		"- Performance optimization recommendations\n"+
		"- Security best practices\n"+
		"- Troubleshooting common issues\n\n"+
		"Would you like me to elaborate on any of these areas?", query)
}

func (g *GeneralAgent) generateSuggestions(query string) []string {
	return []string{
		"Check system resource usage",
		"Review log files for errors",
		"Verify network connectivity",
		"Update security configurations",
	}
}

func (g *GeneralAgent) generateActions(query string) []Action {
	return []Action{
		{
			Type:        "command",
			Description: "Check system status",
			Command:     "allora monitor status",
			Parameters:  map[string]interface{}{"format": "table"},
			Risk:        "low",
		},
		{
			Type:        "analysis",
			Description: "Analyze system logs",
			Command:     "allora analyze logs --time 24h",
			Parameters:  map[string]interface{}{"pattern": "error|warning"},
			Risk:        "low",
		},
	}
}

func (a *AWSAgent) generateAWSResponse(query string) string {
	return fmt.Sprintf("AWS Analysis for '%s':\n\n"+
		"I can help you with AWS-specific operations including:\n"+
		"- EC2 instance optimization\n"+
		"- S3 bucket security and cost management\n"+
		"- RDS performance tuning\n"+
		"- Lambda function monitoring\n\n"+
		"Please specify which AWS service you'd like to focus on.", query)
}

func (a *AWSAgent) generateAWSSuggestions(query string) []string {
	return []string{
		"Check EC2 instance utilization",
		"Review S3 bucket policies",
		"Monitor RDS performance metrics",
		"Optimize Lambda function costs",
	}
}

func (a *AWSAgent) generateAWSActions(query string) []Action {
	return []Action{
		{
			Type:        "aws-command",
			Description: "List EC2 instances",
			Command:     "aws ec2 describe-instances",
			Parameters:  map[string]interface{}{"region": "us-west-2"},
			Risk:        "low",
		},
	}
}

func (az *AzureAgent) generateAzureResponse(query string) string {
	return fmt.Sprintf("Azure Analysis for '%s':\n\n"+
		"I can assist with Azure operations including:\n"+
		"- Virtual machine management\n"+
		"- Storage account optimization\n"+
		"- Azure SQL performance\n"+
		"- Function app monitoring\n\n"+
		"Which Azure service would you like to explore?", query)
}

func (az *AzureAgent) generateAzureSuggestions(query string) []string {
	return []string{
		"Check VM resource usage",
		"Review storage account access",
		"Monitor Azure SQL performance",
		"Optimize Function app costs",
	}
}

func (az *AzureAgent) generateAzureActions(query string) []Action {
	return []Action{
		{
			Type:        "azure-command",
			Description: "List virtual machines",
			Command:     "az vm list",
			Parameters:  map[string]interface{}{"output": "table"},
			Risk:        "low",
		},
	}
}

func (g *GCPAgent) generateGCPResponse(query string) string {
	return fmt.Sprintf("GCP Analysis for '%s':\n\n"+
		"I can help with Google Cloud operations:\n"+
		"- Compute Engine optimization\n"+
		"- Cloud Storage management\n"+
		"- Cloud SQL monitoring\n"+
		"- Cloud Functions analysis\n\n"+
		"Which GCP service interests you?", query)
}

func (g *GCPAgent) generateGCPSuggestions(query string) []string {
	return []string{
		"Check Compute Engine usage",
		"Review Cloud Storage buckets",
		"Monitor Cloud SQL instances",
		"Optimize Cloud Functions",
	}
}

func (g *GCPAgent) generateGCPActions(query string) []Action {
	return []Action{
		{
			Type:        "gcp-command",
			Description: "List compute instances",
			Command:     "gcloud compute instances list",
			Parameters:  map[string]interface{}{"format": "table"},
			Risk:        "low",
		},
	}
}

func (k *KubernetesAgent) generateK8sResponse(query string) string {
	return fmt.Sprintf("Kubernetes Analysis for '%s':\n\n"+
		"I can assist with Kubernetes operations:\n"+
		"- Pod and container management\n"+
		"- Service mesh configuration\n"+
		"- Resource optimization\n"+
		"- Security policies\n\n"+
		"What aspect of Kubernetes would you like to explore?", query)
}

func (k *KubernetesAgent) generateK8sSuggestions(query string) []string {
	return []string{
		"Check pod resource usage",
		"Review service configurations",
		"Monitor cluster health",
		"Optimize resource requests",
	}
}

func (k *KubernetesAgent) generateK8sActions(query string) []Action {
	return []Action{
		{
			Type:        "kubectl-command",
			Description: "List pods",
			Command:     "kubectl get pods",
			Parameters:  map[string]interface{}{"all-namespaces": "true"},
			Risk:        "low",
		},
	}
}

func (m *MonitoringAgent) generateMonitoringResponse(query string) string {
	return fmt.Sprintf("Monitoring Analysis for '%s':\n\n"+
		"I can help with observability and monitoring:\n"+
		"- Metrics collection and analysis\n"+
		"- Alert configuration\n"+
		"- Dashboard creation\n"+
		"- SLO/SLI monitoring\n\n"+
		"Which monitoring aspect would you like to focus on?", query)
}

func (m *MonitoringAgent) generateMonitoringSuggestions(query string) []string {
	return []string{
		"Set up key performance indicators",
		"Configure alerting rules",
		"Create monitoring dashboards",
		"Implement SLO tracking",
	}
}

func (m *MonitoringAgent) generateMonitoringActions(query string) []Action {
	return []Action{
		{
			Type:        "monitoring-command",
			Description: "Check system metrics",
			Command:     "allora monitor metrics",
			Parameters:  map[string]interface{}{"duration": "1h"},
			Risk:        "low",
		},
	}
}

// ProcessQuery processes a query using available agents
func (m *AgentManager) ProcessQuery(ctx context.Context, queryText string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	// If no agents are available, return a helpful message
	if len(m.agents) == 0 {
		return `I'm AlloraAi, your AI-powered infrastructure assistant! 

I can help you with:
🔧 Cloud infrastructure management (AWS, Azure, GCP)
🚀 Application deployment and scaling
📊 Monitoring and alerting setup
🐛 Troubleshooting and debugging
🔒 Security analysis and compliance
📈 Performance optimization

To get started, you'll need to configure your cloud providers using:
- allora config set
- allora init

For now, I'm running in demo mode. How can I help you today?`, nil
	}
	
	// Create a query object
	query := &Query{
		Text:    queryText,
		Context: make(map[string]interface{}),
	}
	
	// Try to find the best agent for the query
	// For now, just use the first available agent
	for _, agent := range m.agents {
		if agent.IsHealthy() {
			response, err := agent.Query(ctx, query)
			if err != nil {
				continue // Try next agent
			}
			return response.Content, nil
		}
	}
	
	// If no healthy agents, return a fallback response
	return fmt.Sprintf(`I understand you're asking about: "%s"

While I'm currently in demo mode, I can help you with infrastructure management tasks like:
- Setting up monitoring for your applications
- Deploying services to cloud platforms
- Troubleshooting performance issues
- Configuring security policies
- Optimizing resource usage

To enable full AI capabilities, please configure your API keys using:
allora config set openai.api_key YOUR_API_KEY

Would you like me to help you get started with the setup?`, queryText), nil
}
