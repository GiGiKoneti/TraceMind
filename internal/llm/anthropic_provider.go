package llm

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicProvider struct {
	client anthropic.Client
	model  anthropic.Model
}

func NewAnthropicProvider(model, apiKey string) (*AnthropicProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("anthropic api key is required")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	return &AnthropicProvider{
		client: client,
		model:  anthropic.Model(model),
	}, nil
}

func (p *AnthropicProvider) Generate(ctx context.Context, prompt string) (string, error) {
	message, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     p.model,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		return "", fmt.Errorf("anthropic generation failed: %w", err)
	}

	if len(message.Content) == 0 {
		return "", fmt.Errorf("no response from anthropic")
	}

	return message.Content[0].Text, nil
}

func (p *AnthropicProvider) GenerateStream(ctx context.Context, prompt string, onToken func(string)) error {
	stream := p.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     p.model,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	for stream.Next() {
		event := stream.Current()

		if event.Type == "content_block_delta" {
			if event.Delta.Type == "text_delta" {
				onToken(event.Delta.Text)
			}
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("anthropic stream error: %w", err)
	}

	return nil
}
