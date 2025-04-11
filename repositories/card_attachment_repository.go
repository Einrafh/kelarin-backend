package repositories

import (
	"kelarin-backend/database"
	"kelarin-backend/models"
)

// CreateCardAttachment creates a new card attachment.
func CreateCardAttachment(attachment *models.CardAttachment) error {
	return database.DB.Create(attachment).Error
}

// GetAttachmentsByCardID retrieves all card attachments for a given card.
func GetAttachmentsByCardID(cardID uint, attachments *[]models.CardAttachment) error {
	return database.DB.Where("card_id = ?", cardID).Find(attachments).Error
}

// GetCardAttachmentByID retrieves a card attachment by its ID.
func GetCardAttachmentByID(id uint, attachment *models.CardAttachment) error {
	return database.DB.First(attachment, id).Error
}

// UpdateCardAttachment updates an existing card attachment.
func UpdateCardAttachment(attachment *models.CardAttachment) error {
	return database.DB.Save(attachment).Error
}

// DeleteCardAttachment deletes a card attachment by its ID.
func DeleteCardAttachment(id uint) error {
	return database.DB.Delete(&models.CardAttachment{}, id).Error
}
