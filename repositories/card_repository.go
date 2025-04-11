package repositories

import (
	"kelarin-backend/database"
	"kelarin-backend/models"
)

// === Card Functions ===

// CreateCard creates a new card.
func CreateCard(card *models.Card) error {
	return database.DB.Create(card).Error
}

// GetCardsByListID retrieves cards for a given list.
func GetCardsByListID(listID uint, cards *[]models.Card) error {
	return database.DB.Where("list_id = ?", listID).
		Preload("Subtasks").
		Preload("Assignees.User").
		Preload("Attachments").
		Preload("Labels").
		Preload("Comments.User").
		Find(cards).Error
}

// GetCardByID retrieves a card by its ID, preloading its associations.
func GetCardByID(id uint, card *models.Card) error {
	return database.DB.
		Preload("Subtasks").
		Preload("Assignees.User").
		Preload("Attachments").
		Preload("Labels").
		Preload("Comments.User").
		First(card, id).Error
}

// UpdateCard updates an existing card.
func UpdateCard(card *models.Card) error {
	return database.DB.Save(card).Error
}

// DeleteCard deletes a card by its ID.
func DeleteCard(id uint) error {
	return database.DB.Delete(&models.Card{}, id).Error
}
