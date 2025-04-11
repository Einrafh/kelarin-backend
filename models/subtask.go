package models

// Subtask represents a subtask or detail task attached to a Card.
type Subtask struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `gorm:"not null" json:"title"`
	IsDone bool   `gorm:"default:false" json:"is_done"`
	CardID uint   `gorm:"not null" json:"card_id"`

	// The card this subtask belongs to.
	Card Card `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE" json:"-"`
}
