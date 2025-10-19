package domain

import "time"

type CVEFinding struct {
	ID           int       `json:"id"`
	ImageID      int       `json:"image_id"`
	CVEID        string    `json:"cve_id"`
	Severity     string    `json:"severity"`
	Description  string    `json:"description"`
	FixAvailable bool      `json:"fix_available"`
	CreatedAt    time.Time `json:"created_at"`
}

