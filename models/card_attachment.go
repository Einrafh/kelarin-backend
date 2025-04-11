package models

import "time"

// CardAttachment represents an attachment for a card, which can be a file or a link.
type CardAttachment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CardID    uint      `gorm:"not null" json:"card_id"`
	URL       string    `json:"url"`       // URL or file path
	FileName  string    `json:"file_name"` // Optional: file name of the uploaded file
	CreatedAt time.Time `json:"created_at"`

	Card Card `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE" json:"-"`
}
