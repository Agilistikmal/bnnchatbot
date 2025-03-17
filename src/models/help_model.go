package models

import "time"

type Help struct {
	JID         string    `gorm:"primaryKey" json:"jid,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Name        string    `json:"name,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	DisplayTime string    `json:"display_time,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
