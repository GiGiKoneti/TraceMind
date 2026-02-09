package memory

import (
	"sync"
	"time"

	"github.com/gigikoneti/tracemind/internal/models"
)

// Store represents the symbolic memory of the system.
type Store struct {
	mu           sync.RWMutex
	recentTraces []models.Trace
	maxTraces    int
}

func NewStore(maxItems int) *Store {
	return &Store{
		recentTraces: make([]models.Trace, 0),
		maxTraces:    maxItems,
	}
}

// AddTrace adds a trace to memory and rotates if necessary.
func (s *Store) AddTrace(trace models.Trace) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.recentTraces = append(s.recentTraces, trace)
	if len(s.recentTraces) > s.maxTraces {
		s.recentTraces = s.recentTraces[1:]
	}
}

// GetHealth computes an aggregate health view from memory.
func (s *Store) GetHealth() models.SystemHealth {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var totalSpans int
	var errorSpans int
	serviceLatencies := make(map[string][]float64)

	for _, t := range s.recentTraces {
		for _, span := range t.Spans {
			totalSpans++
			if span.Status.Code == "ERROR" {
				errorSpans++
			}
			serviceLatencies[span.Name] = append(serviceLatencies[span.Name], span.LatencyMs())
		}
	}

	health := models.SystemHealth{
		LastUpdate: time.Now(),
	}

	if totalSpans > 0 {
		health.RecentErrorRate = float64(errorSpans) / float64(totalSpans)
	}

	// Simple heuristic for slowest services
	for svc, latencies := range serviceLatencies {
		var avg float64
		for _, l := range latencies {
			avg += l
		}
		avg /= float64(len(latencies))
		if avg > 500 {
			health.SlowestServices = append(health.SlowestServices, svc)
		}
	}

	return health
}
