package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type OllamaProvider struct {
	llm   *ollama.LLM
	model string
}

func NewOllamaProvider(endpoint, model string) (*OllamaProvider, error) {
	opts := []ollama.Option{
		ollama.WithModel(model),
	}

	if endpoint != "" {
		opts = append(opts, ollama.WithServerURL(endpoint))
	}

	l, err := ollama.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}

	return &OllamaProvider{
		llm:   l,
		model: model,
	}, nil
}

func (p *OllamaProvider) Generate(ctx context.Context, prompt string) (string, error) {
	completion, err := llms.GenerateFromSinglePrompt(ctx, p.llm, prompt)
	if err != nil {
		return "", fmt.Errorf("ollama generation failed: %w", err)
	}

	return completion, nil
}

func (p *OllamaProvider) GenerateStream(ctx context.Context, prompt string, onToken func(string)) error {
	_, err := llms.GenerateFromSinglePrompt(ctx, p.llm, prompt,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			onToken(string(chunk))
			return nil
		}),
	)

	if err != nil {
		return fmt.Errorf("ollama stream failed: %w", err)
	}

	return nil
}
