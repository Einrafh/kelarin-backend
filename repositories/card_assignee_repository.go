package repositories

import (
	"kelarin-backend/database"
	"kelarin-backend/models"
)

// CreateCardAssignee adds a new card assignee.
func CreateCardAssignee(assignee *models.CardAssignee) error {
	return database.DB.Create(assignee).Error
}

// GetAssigneesByCardID retrieves all card assignees for a given card.
func GetAssigneesByCardID(cardID uint, assignees *[]models.CardAssignee) error {
	return database.DB.
		Where("card_id = ?", cardID).
		Preload("User").
		Find(assignees).Error
}

// GetCardAssignee retrieves a specific card assignee by card ID and user ID.
func GetCardAssignee(cardID, userID uint, assignee *models.CardAssignee) error {
	return database.DB.
		Where("card_id = ? AND user_id = ?", cardID, userID).
		Preload("User").
		First(assignee).Error
}

// DeleteCardAssignee removes a card assignee.
func DeleteCardAssignee(cardID, userID uint) error {
	return database.DB.Delete(&models.CardAssignee{}, "card_id = ? AND user_id = ?", cardID, userID).Error
}
