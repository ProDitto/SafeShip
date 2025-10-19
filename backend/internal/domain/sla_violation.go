package domain

import "time"

type SLAViolation struct {
	ID              int       `json:"id"`
	TenantNamespace string    `json:"tenant_namespace"`
	CVEFindingID    int       `json:"cve_finding_id"`
	Status          string    `json:"status"` // e.g., "active", "resolved"
	CreatedAt       time.Time `json:"created_at"`
	ResolvedAt      time.Time `json:"resolved_at"`
}

