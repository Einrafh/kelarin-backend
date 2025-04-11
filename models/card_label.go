package models

import "time"

// CardLabel represents a label on a card.
type CardLabel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CardID    uint      `gorm:"not null" json:"card_id"`
	Name      string    `gorm:"not null" json:"name"`
	Color     string    `json:"color"` // e.g., hex code or color name
	CreatedAt time.Time `json:"created_at"`

	Card Card `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE" json:"-"`
}
