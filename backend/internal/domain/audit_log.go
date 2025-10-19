package domain

import "time"

type AuditLog struct {
	ID              int                    `json:"id"`
	TenantNamespace string                 `json:"tenant_namespace"`
	Action          string                 `json:"action"`
	Actor           string                 `json:"actor"` // e.g., user ID, API key ID, system process
	Details         map[string]interface{} `json:"details"`
	Timestamp       time.Time              `json:"timestamp"`
}

