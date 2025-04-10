package dto

import "kelarin-backend/models"

// WorkspaceUserResponse represents the collaborator information in the response.
type WorkspaceUserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// UserResponse represents the basic user info in the workspace response.
type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
}

// WorkspaceResponse represents the workspace data to be sent as a response.
type WorkspaceResponse struct {
	ID               uint                    `json:"id"`
	Title            string                  `json:"title"`
	Purpose          string                  `json:"purpose"`
	Description      string                  `json:"description"`
	WorkspacePicture string                  `json:"workspace_picture"`
	WorkspaceBanner  string                  `json:"workspace_banner"`
	Owner            UserResponse            `json:"owner"`
	Collaborators    []WorkspaceUserResponse `json:"collaborators"`
}

// NewWorkspaceResponse converts a Workspace model to a WorkspaceResponse DTO.
func NewWorkspaceResponse(w *models.Workspace) WorkspaceResponse {
	collaborators := make([]WorkspaceUserResponse, len(w.Collaborators))
	for i, c := range w.Collaborators {
		collaborators[i] = WorkspaceUserResponse{
			ID:       c.User.ID,
			FullName: c.User.FullName,
			Email:    c.User.Email,
			Role:     c.Role,
		}
	}

	return WorkspaceResponse{
		ID:               w.ID,
		Title:            w.Title,
		Purpose:          w.Purpose,
		Description:      w.Description,
		WorkspacePicture: w.WorkspacePicture,
		WorkspaceBanner:  w.WorkspaceBanner,
		Owner: UserResponse{
			ID:       w.Owner.ID,
			FullName: w.Owner.FullName,
			Email:    w.Owner.Email,
		},
		Collaborators: collaborators,
	}
}
