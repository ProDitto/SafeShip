package domain

import "time"

type Notification struct {
	ID              int       `json:"id"`
	TenantNamespace string    `json:"tenant_namespace"`
	Type            string    `json:"type"` // e.g., "SLA_VIOLATION", "BUILD_COMPLETE"
	Payload         string    `json:"payload"`
	SentAt          time.Time `json:"sent_at"`
	Status          string    `json:"status"` // e.g., "sent", "failed"
}

