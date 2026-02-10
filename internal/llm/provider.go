package llm

import (
	"context"
	"fmt"

	"github.com/gigikoneti/tracemind/internal/models"
)

type LLMProvider interface {
	Generate(ctx context.Context, prompt string) (string, error)
	GenerateStream(ctx context.Context, prompt string, onToken func(string)) error
}

type ProviderConfig struct {
	Provider AIProvider
	Model    string
	Endpoint string
	APIKey   string
}

type AIProvider string

const (
	ProviderOpenAI    AIProvider = "openai"
	ProviderAnthropic AIProvider = "anthropic"
	ProviderOllama    AIProvider = "ollama"
)

func NewEngineFromConnection(conn models.AIConnection) (*Engine, error) {
	var provider LLMProvider
	var err error

	switch conn.Provider {
	case models.ProviderOpenAI:
		config, configErr := conn.GetOpenAIConfig()
		if configErr != nil {
			return nil, configErr
		}
		provider, err = NewOpenAIProvider(config.APIEndpoint, config.Model, config.APIKey)
		if err != nil {
			return nil, err
		}

	case models.ProviderAnthropic:
		config, configErr := conn.GetAnthropicConfig()
		if configErr != nil {
			return nil, configErr
		}
		provider, err = NewAnthropicProvider(config.Model, config.APIKey)
		if err != nil {
			return nil, err
		}

	case models.ProviderOllama:
		config, configErr := conn.GetOllamaConfig()
		if configErr != nil {
			return nil, configErr
		}
		provider, err = NewOllamaProvider(config.Endpoint, config.Model)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported provider: %s", conn.Provider)
	}

	return &Engine{
		provider: provider,
		config: ProviderConfig{
			Provider: AIProvider(conn.Provider),
			Model:    conn.Config["model"].(string),
		},
	}, nil
}

func (e *Engine) GenerateText(ctx context.Context, prompt string) (string, error) {
	return e.provider.Generate(ctx, prompt)
}

func (e *Engine) GenerateTextStream(ctx context.Context, prompt string, onToken func(string)) error {
	return e.provider.GenerateStream(ctx, prompt, onToken)
}
