package models

import "time"

type InfrastructureDesign struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Components  []Component `json:"components"`
	Metadata    DesignMeta  `json:"metadata"`
}

type Component struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Spec       map[string]interface{} `json:"spec"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type DesignMeta struct {
	GeneratedBy string    `json:"generated_by"`
	Prompt      string    `json:"prompt"`
	Provider    string    `json:"provider"`
	Model       string    `json:"model"`
	GeneratedAt time.Time `json:"generated_at"`
}
