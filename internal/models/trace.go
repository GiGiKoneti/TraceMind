package models

import "time"

// Attribute represents an OTel metadata key-value pair.
type Attribute struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Span represents an OTel-aligned unit of work.
type Span struct {
	SpanID        string      `json:"span_id"`
	TraceID       string      `json:"trace_id"`
	ParentSpanID  string      `json:"parent_span_id,omitempty"`
	Name          string      `json:"name"`
	Kind          string      `json:"kind"`
	StartTime     time.Time   `json:"start_time"`
	EndTime       time.Time   `json:"end_time"`
	Attributes    []Attribute `json:"attributes,omitempty"`
	Status        Status      `json:"status"`
	ResourceNames []string    `json:"resource_names,omitempty"`
}

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

// LatencyMs returns the duration in milliseconds.
func (s *Span) LatencyMs() float64 {
	return float64(s.EndTime.Sub(s.StartTime).Microseconds()) / 1000.0
}

// Trace represents a collection of OTel spans.
type Trace struct {
	TraceID string `json:"trace_id"`
	Spans   []Span `json:"spans"`
}

// SymbolicFact represents a pre-computed insight about the trace.
type SymbolicFact struct {
	Type        string `json:"type"`
	Service     string `json:"service"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// SystemHealth represents global context for symbolic memory.
type SystemHealth struct {
	RecentErrorRate float64   `json:"recent_error_rate"`
	SlowestServices []string  `json:"slowest_services"`
	LastUpdate      time.Time `json:"last_update"`
}

// TraceAnalysis combines raw data with symbolic reasoning and AI streaming output.
type TraceAnalysis struct {
	Trace         Trace          `json:"trace"`
	SymbolicFacts []SymbolicFact `json:"symbolic_facts"`
	SystemContext SystemHealth   `json:"system_context"`
	AIExplanation string         `json:"ai_explanation"` // For unary responses
}
