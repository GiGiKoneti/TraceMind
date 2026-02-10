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

	connectionStore := handlers.NewConnectionStore()
	connectionHandler := &handlers.ConnectionHandler{
		Store: connectionStore,
	}
	designHandler := &handlers.AIDesignHandler{
		ConnectionStore: connectionStore,
	}

	// Common CORS middleware for all routes
	withCORS := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				return
			}
			h(w, r)
		}
	}

	// Trace analysis routes (existing)
	http.HandleFunc("/api/analyze", withCORS(traceHandler.AnalyzeTraceStream))
	http.HandleFunc("/api/evaluate", withCORS(traceHandler.Evaluate))

	// AI Connection routes (new)
	http.HandleFunc("/api/connections", withCORS(connectionHandler.ListConnections))
	http.HandleFunc("/api/connections/create", withCORS(connectionHandler.CreateConnection))
	http.HandleFunc("/api/connections/test", withCORS(connectionHandler.TestConnection))
	http.HandleFunc("/api/connections/delete", withCORS(connectionHandler.DeleteConnection))

	// AI Design generation routes (new)
	http.HandleFunc("/api/design/generate", withCORS(designHandler.GenerateDesign))
	http.HandleFunc("/api/design/generate-stream", withCORS(designHandler.GenerateDesignStream))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("TraceMind AI Adapter starting on :%s (using model: %s)", port, modelName)
	log.Printf("Available endpoints:")
	log.Printf("  - Trace Analysis: /api/analyze, /api/evaluate")
	log.Printf("  - AI Connections: /api/connections, /api/connections/create, /api/connections/test, /api/connections/delete")
	log.Printf("  - Design Generation: /api/design/generate, /api/design/generate-stream")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
