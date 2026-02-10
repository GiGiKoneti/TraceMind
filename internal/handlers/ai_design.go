package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gigikoneti/tracemind/internal/llm"
	"github.com/gigikoneti/tracemind/internal/models"
	"github.com/google/uuid"
)

type AIDesignHandler struct {
	ConnectionStore *ConnectionStore
}

func (h *AIDesignHandler) GenerateDesign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Prompt       string `json:"prompt"`
		ConnectionID string `json:"connection_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	if req.ConnectionID == "" {
		http.Error(w, "Connection ID is required", http.StatusBadRequest)
		return
	}

	conn, ok := h.ConnectionStore.Get(req.ConnectionID)
	if !ok {
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}

	engine, err := llm.NewEngineFromConnection(conn)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize LLM engine: %v", err), http.StatusInternalServerError)
		return
	}

	systemPrompt := buildInfrastructurePrompt(req.Prompt)

	response, err := engine.GenerateText(r.Context(), systemPrompt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Design generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	design, err := parseDesignResponse(response, req.Prompt, conn)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse design: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(design)
}

func (h *AIDesignHandler) GenerateDesignStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Prompt       string `json:"prompt"`
		ConnectionID string `json:"connection_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	if req.ConnectionID == "" {
		http.Error(w, "Connection ID is required", http.StatusBadRequest)
		return
	}

	conn, ok := h.ConnectionStore.Get(req.ConnectionID)
	if !ok {
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}

	engine, err := llm.NewEngineFromConnection(conn)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize LLM engine: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	systemPrompt := buildInfrastructurePrompt(req.Prompt)

	var fullResponse strings.Builder

	err = engine.GenerateTextStream(r.Context(), systemPrompt, func(token string) {
		fullResponse.WriteString(token)
		fmt.Fprintf(w, "event: token\ndata: %s\n\n", token)
		w.(http.Flusher).Flush()
	})

	if err != nil {
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
		return
	}

	design, err := parseDesignResponse(fullResponse.String(), req.Prompt, conn)
	if err != nil {
		fmt.Fprintf(w, "event: error\ndata: Failed to parse design: %s\n\n", err.Error())
		return
	}

	designJSON, _ := json.Marshal(design)
	fmt.Fprintf(w, "event: design\ndata: %s\n\n", designJSON)
	fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
}

func buildInfrastructurePrompt(userPrompt string) string {
	return fmt.Sprintf(`You are an expert infrastructure architect specializing in Kubernetes and cloud-native systems. Generate a complete infrastructure design in JSON format based on the user's requirements.

User Request: %s

Output Format (JSON):
{
  "name": "design-name",
  "description": "brief description of the infrastructure",
  "version": "1.0.0",
  "components": [
    {
      "id": "unique-id",
      "name": "component-name",
      "type": "Kubernetes",
      "apiVersion": "apps/v1",
      "kind": "Deployment",
      "spec": {
        "replicas": 3,
        "selector": {...},
        "template": {...}
      },
      "metadata": {
        "namespace": "default",
        "labels": {...}
      }
    }
  ]
}

Rules:
1. Include ALL necessary components (Deployments, Services, ConfigMaps, PersistentVolumeClaims, etc.)
2. Use valid Kubernetes API versions (apps/v1, v1, networking.k8s.io/v1, etc.)
3. Follow best practices:
   - High availability (multiple replicas)
   - Resource limits and requests
   - Health checks (liveness/readiness probes)
   - Security (non-root users, read-only filesystems where appropriate)
4. For monitoring: Include Prometheus ServiceMonitor if requested
5. For databases: Include StatefulSets with persistent storage
6. Return ONLY valid JSON, no markdown code blocks, no explanations

Generate the complete infrastructure design now:`, userPrompt)
}

func parseDesignResponse(response, prompt string, conn models.AIConnection) (*models.InfrastructureDesign, error) {
	response = strings.TrimSpace(response)

	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var design models.InfrastructureDesign
	if err := json.Unmarshal([]byte(response), &design); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if design.Name == "" {
		design.Name = fmt.Sprintf("design-%s", uuid.New().String()[:8])
	}
	if design.Version == "" {
		design.Version = "1.0.0"
	}

	design.Metadata = models.DesignMeta{
		GeneratedBy: "TraceMind AI Adapter",
		Prompt:      prompt,
		Provider:    string(conn.Provider),
		Model:       getModelFromConfig(conn),
		GeneratedAt: time.Now(),
	}

	for i := range design.Components {
		if design.Components[i].ID == "" {
			design.Components[i].ID = uuid.New().String()
		}
	}

	return &design, nil
}

func getModelFromConfig(conn models.AIConnection) string {
	if model, ok := conn.Config["model"].(string); ok {
		return model
	}
	return "unknown"
}

func (h *AIDesignHandler) ValidateDesign(ctx context.Context, design *models.InfrastructureDesign) error {
	if design.Name == "" {
		return fmt.Errorf("design name is required")
	}

	if len(design.Components) == 0 {
		return fmt.Errorf("design must have at least one component")
	}

	for i, comp := range design.Components {
		if comp.Name == "" {
			return fmt.Errorf("component %d: name is required", i)
		}
		if comp.Kind == "" {
			return fmt.Errorf("component %s: kind is required", comp.Name)
		}
		if comp.APIVersion == "" {
			return fmt.Errorf("component %s: apiVersion is required", comp.Name)
		}
	}

	return nil
}
