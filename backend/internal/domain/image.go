package domain

import (
	"time"
)

type Image struct {
	ID              int       `json:"id"`
	TenantNamespace string    `json:"tenant_namespace"`
	Digest          string    `json:"digest"`
	Tags            []string  `json:"tags"`
	SLSALevel       int       `json:"slsa_level"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

