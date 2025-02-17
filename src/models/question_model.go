package models

import (
	"time"

	"github.com/lib/pq"
)

type Question struct {
	ID        string         `gorm:"primaryKey" json:"id,omitempty"`
	Triggers  pq.StringArray `gorm:"type:text[]" json:"triggers,omitempty"`
	Responses pq.StringArray `gorm:"type:text[]" json:"responses,omitempty"`
	CreatedAt *time.Time     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}
