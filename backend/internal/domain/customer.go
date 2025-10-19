package domain

import "time"

type Customer struct {
	Namespace    string    `json:"namespace"`
	Name         string    `json:"name"`
	ContactInfo  string    `json:"contact_info"`
	SLATier      string    `json:"sla_tier"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

