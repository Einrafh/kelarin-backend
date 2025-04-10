package dto

import "kelarin-backend/models"

// ProfileResponse mengembalikan informasi profil user beserta workspace.
type ProfileResponse struct {
	ID               uint                `json:"id"`
	FullName         string              `json:"fullname"`
	Email            string              `json:"email"`
	UserType         string              `json:"user_type"`
	OwnedWorkspaces  []WorkspaceResponse `json:"owned_workspaces"`
	CollabWorkspaces []WorkspaceResponse `json:"collab_workspaces"`
}

// NewProfileResponse mengubah model User menjadi ProfileResponse.
// Untuk collabWorkspaces, kita harus mengekstrak Workspace dari join table WorkspaceUser.
func NewProfileResponse(user *models.User) ProfileResponse {
	// Mapping owned workspaces menggunakan helper yang sudah ada.
	owned := make([]WorkspaceResponse, len(user.OwnedWorkspaces))
	for i, w := range user.OwnedWorkspaces {
		owned[i] = NewWorkspaceResponse(&w)
	}

	// Mapping collab workspaces: pastikan pada preload di query sudah memuat relasi WorkspaceUser.Workspace, misalnya dengan Preload("CollabWorkspaces.Workspace")
	collab := make([]WorkspaceResponse, 0)
	for _, wUser := range user.CollabWorkspaces {
		// Pastikan wUser.Workspace sudah terisi
		if wUser.Workspace.ID != 0 {
			collab = append(collab, NewWorkspaceResponse(&wUser.Workspace))
		}
	}

	return ProfileResponse{
		ID:               user.ID,
		FullName:         user.FullName,
		Email:            user.Email,
		UserType:         user.UserType,
		OwnedWorkspaces:  owned,
		CollabWorkspaces: collab,
	}
}
