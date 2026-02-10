package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gigikoneti/tracemind/internal/llm"
	"github.com/gigikoneti/tracemind/internal/models"
	"github.com/google/uuid"
)

type ConnectionStore struct {
	connections map[string]models.AIConnection
	mu          sync.RWMutex
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		connections: make(map[string]models.AIConnection),
	}
}

func (s *ConnectionStore) Add(conn models.AIConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[conn.ID] = conn
}

func (s *ConnectionStore) Get(id string) (models.AIConnection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conn, ok := s.connections[id]
	return conn, ok
}

func (s *ConnectionStore) List() []models.AIConnection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conns := make([]models.AIConnection, 0, len(s.connections))
	for _, conn := range s.connections {
		conns = append(conns, conn)
	}
	return conns
}

func (s *ConnectionStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.connections[id]; exists {
		delete(s.connections, id)
		return true
	}
	return false
}

type ConnectionHandler struct {
	Store *ConnectionStore
}

func (h *ConnectionHandler) CreateConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var conn models.AIConnection
	if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := conn.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	conn.ID = uuid.New().String()
	conn.Status = "untested"

	h.Store.Add(conn)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conn)
}

func (h *ConnectionHandler) ListConnections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	connections := h.Store.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connections)
}

func (h *ConnectionHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ConnectionID string `json:"connection_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	conn, ok := h.Store.Get(req.ConnectionID)
	if !ok {
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}

	engine, err := llm.NewEngineFromConnection(conn)
	if err != nil {
		conn.Status = "error"
		h.Store.Add(conn)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	testPrompt := "Reply with 'OK' if you can read this."
	response, err := engine.GenerateText(r.Context(), testPrompt)
	if err != nil {
		conn.Status = "error"
		h.Store.Add(conn)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	conn.Status = "connected"
	h.Store.Add(conn)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "connected",
		"message":  "Connection successful",
		"response": response,
	})
}

func (h *ConnectionHandler) DeleteConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ConnectionID string `json:"connection_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !h.Store.Delete(req.ConnectionID) {
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Connection deleted",
	})
}
