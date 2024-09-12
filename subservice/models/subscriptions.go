package models

import "time"

type Subscriptions struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Link       string    `gorm:"not null" json:"link"`
	Data       string    `gorm:"not null" json:"-"`
	Active     bool      `gorm:"default:false" json:"active"`
	AutoUpdate bool      `gorm:"default:false" json:"auto_update" form:"auto_update"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
