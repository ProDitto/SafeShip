package domain

import "time"

type SBOMRecord struct {
	ID        int       `json:"id"`
	ImageID   int       `json:"image_id"`
	Format    string    `json:"format"`
	URI       string    `json:"uri"`
	CreatedAt time.Time `json:"created_at"`
}

