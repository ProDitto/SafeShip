package domain

import "time"

type CustomerImageUsage struct {
	ID              int       `json:"id"`
	TenantNamespace string    `json:"tenant_namespace"`
	ImageID         int       `json:"image_id"`
	VersionPinned   bool      `json:"version_pinned"`
	RuntimeInfo     string    `json:"runtime_info"` // e.g., cluster name, environment
	CreatedAt       time.Time `json:"created_at"`
}

