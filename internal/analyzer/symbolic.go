package analyzer

import (
	"fmt"

	"github.com/gigikoneti/tracemind/internal/models"
)

// AnalyzeTrace performs symbolic reasoning on a trace to extract key facts.
func AnalyzeTrace(trace models.Trace) []models.SymbolicFact {
	var facts []models.SymbolicFact

	if len(trace.Spans) == 0 {
		return facts
	}

	var maxLatencySpan models.Span
	first := true
	for _, span := range trace.Spans {
		latency := span.LatencyMs()
		if first || latency > maxLatencySpan.LatencyMs() {
			maxLatencySpan = span
			first = false
		}
	}

	maxLat := maxLatencySpan.LatencyMs()
	if maxLat > 800 {
		facts = append(facts, models.SymbolicFact{
			Type:        "LATENCY_BOTTLENECK",
			Service:     maxLatencySpan.Name,
			Description: fmt.Sprintf("Service '%s' is a bottleneck with %.2fms latency.", maxLatencySpan.Name, maxLat),
			Severity:    "critical",
		})
	} else if maxLat > 400 {
		facts = append(facts, models.SymbolicFact{
			Type:        "LATENCY_WARNING",
			Service:     maxLatencySpan.Name,
			Description: fmt.Sprintf("Service '%s' has elevated latency: %.2fms.", maxLatencySpan.Name, maxLat),
			Severity:    "warning",
		})
	}

	for _, span := range trace.Spans {
		if span.Status.Code == "ERROR" {
			isOrigin := true
			if span.ParentSpanID != "" {
				for _, p := range trace.Spans {
					if p.SpanID == span.ParentSpanID && p.Status.Code == "ERROR" {
						isOrigin = false
						break
					}
				}
			}

			if isOrigin {
				facts = append(facts, models.SymbolicFact{
					Type:        "ERROR_ORIGIN",
					Service:     span.Name,
					Description: fmt.Sprintf("Error originated in service '%s': %s", span.Name, span.Status.Message),
					Severity:    "critical",
				})
			}
		}
	}

	return facts
}
