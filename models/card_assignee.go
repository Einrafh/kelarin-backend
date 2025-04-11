package models

// CardAssignee links a Card with an assigned User.
type CardAssignee struct {
	CardID uint `gorm:"primaryKey" json:"card_id"`
	UserID uint `gorm:"primaryKey" json:"user_id"`

	// The user assigned to the card.
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}
