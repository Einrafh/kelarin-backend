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

// GetWorkspaceIDByListID retrieves the workspace_id associated with a board list.
func GetWorkspaceIDByListID(listID uint) (uint, error) {
	var boardList models.BoardList
	if err := database.DB.Select("workspace_id").First(&boardList, listID).Error; err != nil {
		return 0, err
	}
	return boardList.WorkspaceID, nil
}

// UpdateBoardList updates an existing board list.
func UpdateBoardList(list *models.BoardList) error {
	return database.DB.Save(list).Error
}

// DeleteBoardList deletes a board list by its ID.
func DeleteBoardList(id uint) error {
	return database.DB.Delete(&models.BoardList{}, id).Error
}
