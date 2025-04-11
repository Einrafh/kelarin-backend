package repositories

import (
	"kelarin-backend/database"
	"kelarin-backend/models"
)

// CreateSubtask creates a new subtask for a card.
func CreateSubtask(subtask *models.Subtask) error {
	return database.DB.Create(subtask).Error
}

// GetSubtasksByCard retrieves subtasks for a given card.
func GetSubtasksByCard(cardID uint, subtasks *[]models.Subtask) error {
	return database.DB.Where("card_id = ?", cardID).Find(subtasks).Error
}

// GetSubtaskByID retrieves a subtask by its ID.
func GetSubtaskByID(id uint, subtask *models.Subtask) error {
	return database.DB.First(subtask, id).Error
}

// UpdateSubtask updates an existing subtask.
func UpdateSubtask(subtask *models.Subtask) error {
	return database.DB.Save(subtask).Error
}

// DeleteSubtask deletes a subtask by its ID.
func DeleteSubtask(id uint) error {
	return database.DB.Delete(&models.Subtask{}, id).Error
}
