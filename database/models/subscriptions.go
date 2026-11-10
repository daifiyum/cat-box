package models

import (
	"time"
)

type Subscriptions struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"not null" json:"name"`
	Link       string `gorm:"not null" json:"link"`
	Data       string `gorm:"not null" json:"-"`
	Active     bool   `gorm:"default:false" json:"active"`
	AutoUpdate bool   `gorm:"default:false" json:"auto_update" form:"auto_update"`
	Insecure   bool   `gorm:"not null;default:false" json:"insecure" form:"insecure"`
	// UserAgent   string    `gorm:"not null" json:"user_agent"`
	SortOrder   int       `gorm:"not null" json:"sort_order"`
	UpdatedTime time.Time `gorm:"not null" json:"updated_at"`
}
