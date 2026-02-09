package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gigikoneti/tracemind/internal/analyzer"
	"github.com/gigikoneti/tracemind/internal/llm"
	"github.com/gigikoneti/tracemind/internal/memory"
	"github.com/gigikoneti/tracemind/internal/models"
)

type TraceHandler struct {
	Engine *llm.Engine
	Memory *memory.Store
}

func (h *TraceHandler) AnalyzeTraceStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	useStructured := r.URL.Query().Get("structured") != "false"

	var trace models.Trace
	if err := json.NewDecoder(r.Body).Decode(&trace); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.Memory.AddTrace(trace)
	facts := analyzer.AnalyzeTrace(trace)
	health := h.Memory.GetHealth()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Send initial metadata
	initialData := map[string]interface{}{
		"facts":  facts,
		"health": health,
		"trace":  trace,
	}
	initialJSON, _ := json.Marshal(initialData)
	fmt.Fprintf(w, "event: metadata\ndata: %s\n\n", initialJSON)
	w.(http.Flusher).Flush()

	err := h.Engine.ExplainTraceStream(r.Context(), trace, facts, health, useStructured, func(token string) {
		fmt.Fprintf(w, "event: token\ndata: %s\n\n", token)
		w.(http.Flusher).Flush()
	})

	if err != nil {
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
	} else {
		fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
	}
}

func (h *TraceHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Trace       models.Trace          `json:"trace"`
		Facts       []models.SymbolicFact `json:"facts"`
		Explanation string                `json:"explanation"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	score, err := h.Engine.EvaluateExplanation(r.Context(), req.Trace, req.Facts, req.Explanation)
	if err != nil {
		http.Error(w, "Evaluation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]string{"evaluation": score})
}
