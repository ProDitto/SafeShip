package domain

import "time"

type BuildEvent struct {
	ID              int       `json:"id"`
	TenantNamespace string    `json:"tenant_namespace"`
	ImageID         *int      `json:"image_id,omitempty"` // Pointer to allow null
	TriggerType     string    `json:"trigger_type"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

