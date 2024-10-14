package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscriptions struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Link        string    `gorm:"not null" json:"link"`
	Data        string    `gorm:"not null" json:"-"`
	Active      bool      `gorm:"default:false" json:"active"`
	AutoUpdate  bool      `gorm:"default:false" json:"auto_update" form:"auto_update"`
	SortOrder   int       `gorm:"not null" json:"sort_order"`
	UpdatedTime time.Time `gorm:"not null" json:"updated_at"`
}

func (s *Subscriptions) BeforeUpdate(tx *gorm.DB) (err error) {
	var existing Subscriptions
	if err := tx.First(&existing, s.ID).Error; err != nil {
		return err
	}

	if existing.Data != s.Data {
		s.UpdatedTime = time.Now()
	}

	return nil
}
