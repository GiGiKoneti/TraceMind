package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/gigikoneti/tracemind/internal/models"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Engine struct {
	llm      *ollama.LLM
	provider LLMProvider
	config   ProviderConfig
}

func NewEngine(modelName string) (*Engine, error) {
	l, err := ollama.New(ollama.WithModel(modelName))
	if err != nil {
		return nil, err
	}

	ollamaProvider, err := NewOllamaProvider("", modelName)
	if err != nil {
		return nil, err
	}

	return &Engine{
		llm:      l,
		provider: ollamaProvider,
		config: ProviderConfig{
			Provider: ProviderOllama,
			Model:    modelName,
		},
	}, nil
}

// ExplainTraceStream uses the LLM to provide a streaming causal explanation.
func (e *Engine) ExplainTraceStream(ctx context.Context, trace models.Trace, facts []models.SymbolicFact, health models.SystemHealth, useStructured bool, onToken func(string)) error {
	var prompt string
	if useStructured {
		prompt = buildStructuredPrompt(trace, facts, health)
	} else {
		prompt = buildRawPrompt(trace)
	}

	_, err := llms.GenerateFromSinglePrompt(ctx, e.llm, prompt,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			onToken(string(chunk))
			return nil
		}),
	)
	return err
}

// EvaluateExplanation uses LLM-as-a-Judge to score an explanation.
func (e *Engine) EvaluateExplanation(ctx context.Context, trace models.Trace, facts []models.SymbolicFact, explanation string) (string, error) {
	prompt := fmt.Sprintf(`
You are a Senior SRE Auditor. Evaluate the following AI-generated incident explanation based on technical correctness and causal logic.

Trace Data (Simplified):
%v

Computed Symbolic Facts:
%v

AI Explanation to Evaluate:
"""
%s
"""

Task:
Score the explanation from 1-10 on 'Causal Correctness'.
Explain why you gave that score.
Check if the explanation identified the root cause mentioned in the symbolic facts.

Response format:
Score: [1-10]
Rationale: [Brief explanation]
Root Cause Found: [Yes/No]
`, trace, facts, explanation)

	completion, err := llms.GenerateFromSinglePrompt(ctx, e.llm, prompt)
	return completion, err
}

func buildRawPrompt(trace models.Trace) string {
	var sb strings.Builder
	sb.WriteString("Analyze this OTel trace and explain what happened:\n\n")
	for _, s := range trace.Spans {
		sb.WriteString(fmt.Sprintf("- %s: %s [%.2fms]\n", s.Name, s.Status.Code, s.LatencyMs()))
	}
	return sb.String()
}

func buildStructuredPrompt(trace models.Trace, facts []models.SymbolicFact, health models.SystemHealth) string {
	var sb strings.Builder
	sb.WriteString("You are an expert SRE Agent. Analyze this OTel trace using both current telemetry and historical system context.\n\n")

	sb.WriteString("### Global System Context (Symbolic Memory):\n")
	sb.WriteString(fmt.Sprintf("- Overall Error Rate: %.2f%%\n", health.RecentErrorRate*100))
	if len(health.SlowestServices) > 0 {
		sb.WriteString(fmt.Sprintf("- Recent Latency Trends: Services %s have been slow recently.\n", strings.Join(health.SlowestServices, ", ")))
	}

	sb.WriteString("\n### Symbolic Facts for This Trace:\n")
	for _, f := range facts {
		sb.WriteString(fmt.Sprintf("- [%s] %s: %s\n", f.Severity, f.Type, f.Description))
	}

	sb.WriteString("\n### OTel Spans:\n")
	for _, s := range trace.Spans {
		status := s.Status.Code
		if s.Status.Message != "" {
			status += fmt.Sprintf(" (%s)", s.Status.Message)
		}
		sb.WriteString(fmt.Sprintf("- %s: %s [%.2fms]\n", s.Name, status, s.LatencyMs()))
	}

	sb.WriteString("\n### Task:\n")
	sb.WriteString("1. Determine if this is an isolated incident or part of a systemic trend based on the global context.\n")
	sb.WriteString("2. Explain the root cause and propagation.\n")
	sb.WriteString("3. Provide high-priority remediation steps.\n")
	sb.WriteString("\nBe technical, concise, and definitive.")

	return sb.String()
}
