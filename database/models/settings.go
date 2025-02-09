package models

type Setting struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Options     string `json:"options"`
	Description string `json:"description,omitempty"`
}
