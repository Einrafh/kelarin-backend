package repositories

import (
	"kelarin-backend/database"
	"kelarin-backend/models"
)

// CreateCardLabel creates a new card label.
func CreateCardLabel(label *models.CardLabel) error {
	return database.DB.Create(label).Error
}

// GetLabelsByCardID retrieves all card labels for a given card.
func GetLabelsByCardID(cardID uint, labels *[]models.CardLabel) error {
	return database.DB.
		Where("card_id = ?", cardID).
		Find(labels).Error
}

// GetCardLabelByID retrieves a card label by its ID.
func GetCardLabelByID(id uint, label *models.CardLabel) error {
	return database.DB.First(label, id).Error
}

// UpdateCardLabel updates an existing card label.
func UpdateCardLabel(label *models.CardLabel) error {
	return database.DB.Save(label).Error
}

// DeleteCardLabel deletes a card label by its ID.
func DeleteCardLabel(id uint) error {
	return database.DB.Delete(&models.CardLabel{}, id).Error
}
