package repositories

import (
	"kelarin-backend/database"
	"kelarin-backend/models"
)

// CreateBoardList creates a new board list.
func CreateBoardList(list *models.BoardList) error {
	return database.DB.Create(list).Error
}

// GetBoardListsByWorkspace retrieves all board lists for a given workspace.
func GetBoardListsByWorkspace(workspaceID uint, lists *[]models.BoardList) error {
	return database.DB.Where("workspace_id = ?", workspaceID).Preload("Cards").Find(lists).Error
}

// GetBoardListByID retrieves a board list by its ID.
func GetBoardListByID(id uint, list *models.BoardList) error {
	return database.DB.First(list, id).Error
}

// UpdateBoardList updates an existing board list.
func UpdateBoardList(list *models.BoardList) error {
	return database.DB.Save(list).Error
}

// DeleteBoardList deletes a board list by its ID.
func DeleteBoardList(id uint) error {
	return database.DB.Delete(&models.BoardList{}, id).Error
}
