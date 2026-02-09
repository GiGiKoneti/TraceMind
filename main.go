package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gigikoneti/tracemind/internal/handlers"
	"github.com/gigikoneti/tracemind/internal/llm"
	"github.com/gigikoneti/tracemind/internal/memory"
)

func main() {
	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "mistral"
	}

	engine, err := llm.NewEngine(modelName)
	if err != nil {
		log.Fatalf("Failed to initialize LLM engine: %v", err)
	}

	store := memory.NewStore(50)

	traceHandler := &handlers.TraceHandler{
		Engine: engine,
		Memory: store,
	}

	// Common CORS middleware for all routes
	withCORS := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				return
			}
			h(w, r)
		}
	}

	http.HandleFunc("/api/analyze", withCORS(traceHandler.AnalyzeTraceStream))
	http.HandleFunc("/api/evaluate", withCORS(traceHandler.Evaluate))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("TraceMind Research Backend starting on :%s (using model: %s)", port, modelName)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
