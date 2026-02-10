package models

import (
	"fmt"
	"time"
)

type AIProvider string

const (
	ProviderOpenAI    AIProvider = "openai"
	ProviderAnthropic AIProvider = "anthropic"
	ProviderOllama    AIProvider = "ollama"
)

type AIConnection struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Provider    AIProvider             `json:"provider"`
	Config      map[string]interface{} `json:"config"`
	Credentials map[string]string      `json:"credentials,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	Status      string                 `json:"status"`
}

type OpenAIConfig struct {
	APIEndpoint string `json:"api_endpoint"`
	Model       string `json:"model"`
	APIKey      string `json:"api_key"`
}

type AnthropicConfig struct {
	Model  string `json:"model"`
	APIKey string `json:"api_key"`
}

type OllamaConfig struct {
	Endpoint string `json:"endpoint"`
	Model    string `json:"model"`
}

func (c *AIConnection) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("connection name is required")
	}
	if c.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	switch c.Provider {
	case ProviderOpenAI:
		return c.validateOpenAI()
	case ProviderAnthropic:
		return c.validateAnthropic()
	case ProviderOllama:
		return c.validateOllama()
	default:
		return fmt.Errorf("unsupported provider: %s", c.Provider)
	}
}

func (c *AIConnection) validateOpenAI() error {
	model, ok := c.Config["model"].(string)
	if !ok || model == "" {
		return fmt.Errorf("openai model is required")
	}

	apiKey, ok := c.Credentials["api_key"]
	if !ok || apiKey == "" {
		return fmt.Errorf("openai api_key is required")
	}

	return nil
}

func (c *AIConnection) validateAnthropic() error {
	model, ok := c.Config["model"].(string)
	if !ok || model == "" {
		return fmt.Errorf("anthropic model is required")
	}

	apiKey, ok := c.Credentials["api_key"]
	if !ok || apiKey == "" {
		return fmt.Errorf("anthropic api_key is required")
	}

	return nil
}

func (c *AIConnection) validateOllama() error {
	endpoint, ok := c.Config["endpoint"].(string)
	if !ok || endpoint == "" {
		return fmt.Errorf("ollama endpoint is required")
	}

	model, ok := c.Config["model"].(string)
	if !ok || model == "" {
		return fmt.Errorf("ollama model is required")
	}

	return nil
}

func (c *AIConnection) GetOpenAIConfig() (*OpenAIConfig, error) {
	if c.Provider != ProviderOpenAI {
		return nil, fmt.Errorf("connection is not OpenAI provider")
	}

	endpoint, _ := c.Config["api_endpoint"].(string)
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1"
	}

	model, _ := c.Config["model"].(string)
	apiKey, _ := c.Credentials["api_key"]

	return &OpenAIConfig{
		APIEndpoint: endpoint,
		Model:       model,
		APIKey:      apiKey,
	}, nil
}

func (c *AIConnection) GetAnthropicConfig() (*AnthropicConfig, error) {
	if c.Provider != ProviderAnthropic {
		return nil, fmt.Errorf("connection is not Anthropic provider")
	}

	model, _ := c.Config["model"].(string)
	apiKey, _ := c.Credentials["api_key"]

	return &AnthropicConfig{
		Model:  model,
		APIKey: apiKey,
	}, nil
}

func (c *AIConnection) GetOllamaConfig() (*OllamaConfig, error) {
	if c.Provider != ProviderOllama {
		return nil, fmt.Errorf("connection is not Ollama provider")
	}

	endpoint, _ := c.Config["endpoint"].(string)
	model, _ := c.Config["model"].(string)

	return &OllamaConfig{
		Endpoint: endpoint,
		Model:    model,
	}, nil
}
