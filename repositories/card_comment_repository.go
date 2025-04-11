package repositories

import (
	"kelarin-backend/database"
	"kelarin-backend/models"
)

// CreateCardComment creates a new card comment.
func CreateCardComment(comment *models.CardComment) error {
	return database.DB.Create(comment).Error
}

// GetCommentsByCardID retrieves all card comments for a given card.
func GetCommentsByCardID(cardID uint, comments *[]models.CardComment) error {
	return database.DB.
		Where("card_id = ?", cardID).
		Preload("User").
		Find(comments).Error
}

// GetCardCommentByID retrieves a card comment by its ID.
func GetCardCommentByID(id uint, comment *models.CardComment) error {
	return database.DB.Preload("User").First(comment, id).Error
}

// UpdateCardComment updates an existing card comment.
func UpdateCardComment(comment *models.CardComment) error {
	return database.DB.Save(comment).Error
}

// DeleteCardComment deletes a card comment by its ID.
func DeleteCardComment(id uint) error {
	return database.DB.Delete(&models.CardComment{}, id).Error
}
