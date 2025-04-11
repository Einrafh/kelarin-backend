package models

import "time"

// CardComment represents a comment on a card.
type CardComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CardID    uint      `gorm:"not null" json:"card_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Comment   string    `gorm:"not null" json:"comment"`
	CreatedAt time.Time `json:"created_at"`

	Card Card `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE" json:"-"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}
