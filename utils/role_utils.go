package utils

import (
	"errors"

	"kelarin-backend/database"
	"kelarin-backend/models"
)

// CheckRoleInWorkspace retrieves the user's role in a workspace.
// It returns the role as a string (e.g., "owner", "admin", "editor", "viewer").
func CheckRoleInWorkspace(userID, workspaceID uint) (string, error) {
	var workspace models.Workspace
	if err := database.DB.Select("owner_id").First(&workspace, workspaceID).Error; err != nil {
		return "", err
	}

	if workspace.OwnerID == userID {
		return "owner", nil
	}

	var wsUser models.WorkspaceUser
	err := database.DB.
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		First(&wsUser).Error
	if err != nil {
		return "", errors.New("user is not associated with this workspace")
	}

	return wsUser.Role, nil
}

// IsEditorAdminOwner returns true if role is "editor", "admin", or "owner".
func IsEditorAdminOwner(role string) bool {
	return role == "editor" || role == "admin" || role == "owner"
}
