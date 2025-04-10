package repositories

import (
	"kelarin-backend/database"
	"kelarin-backend/models"
)

// CreateWorkspace creates a new workspace in the database.
func CreateWorkspace(workspace *models.Workspace) error {
	return database.DB.Create(workspace).Error
}

// GetAllWorkspacesWithOwner retrieves all workspaces with preloaded Owner and Collaborators.User.
func GetAllWorkspacesWithOwner(workspaces *[]models.Workspace) error {
	return database.DB.
		Preload("Owner").
		Preload("Collaborators.User").
		Find(workspaces).Error
}

// GetWorkspaceByIDWithOwner retrieves a workspace by ID with preloaded Owner and Collaborators.User.
func GetWorkspaceByIDWithOwner(workspaceID string, workspace *models.Workspace) error {
	return database.DB.
		Preload("Owner").
		Preload("Collaborators.User").
		First(workspace, workspaceID).Error
}

// GetWorkspacesAccessibleByUser retrieves the workspaces accessible by a user,
// either as owner or collaborator.
func GetWorkspacesAccessibleByUser(userID uint, workspaces *[]models.Workspace) error {
	subquery := database.DB.
		Table("workspace_users").
		Select("workspace_id").
		Where("user_id = ?", userID)

	return database.DB.
		Preload("Owner").
		Preload("Collaborators.User").
		Where("owner_id = ? OR id IN (?)", userID, subquery).
		Find(workspaces).Error
}

// AddCollaboratorToWorkspaceWithRole adds a collaborator to a workspace with a given role.
// It avoids adding if the user is the owner or is already a collaborator.
func AddCollaboratorToWorkspaceWithRole(workspace *models.Workspace, user *models.User, role string) error {
	if user.ID == workspace.OwnerID {
		return nil
	}

	alreadyExists, err := IsUserAlreadyCollaborator(workspace.ID, user.ID)
	if err != nil {
		return err
	}
	if alreadyExists {
		return nil
	}

	workspaceUser := models.WorkspaceUser{
		UserID:      user.ID,
		WorkspaceID: workspace.ID,
		Role:        role,
	}
	return database.DB.Create(&workspaceUser).Error
}

// AddCollaboratorsByEmails adds multiple collaborators by their emails with the default role "viewer".
// It returns the list of emails that failed to be added.
func AddCollaboratorsByEmails(workspace *models.Workspace, emails []string) ([]string, error) {
	var failedEmails []string
	for _, email := range emails {
		user, err := GetUserByEmail(email)
		if err != nil || user.ID == workspace.OwnerID {
			failedEmails = append(failedEmails, email)
			continue
		}
		if err := AddCollaboratorToWorkspaceWithRole(workspace, user, "viewer"); err != nil {
			failedEmails = append(failedEmails, email)
		}
	}
	return failedEmails, nil
}

// IsUserAlreadyCollaborator checks if a user is already a collaborator in a workspace.
func IsUserAlreadyCollaborator(workspaceID uint, userID uint) (bool, error) {
	var count int64
	err := database.DB.
		Table("workspace_users").
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&count).Error
	return count > 0, err
}

// DeleteWorkspace deletes a workspace by its ID.
func DeleteWorkspace(id string) error {
	return database.DB.Delete(&models.Workspace{}, id).Error
}
