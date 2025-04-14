package dto

import (
	"kelarin-backend/models"
	"kelarin-backend/utils"
)

// ProfileResponse returns user profile information along with workspace details.
type ProfileResponse struct {
	ID               uint                `json:"id"`
	FullName         string              `json:"fullname"`
	Email            string              `json:"email"`
	UserType         string              `json:"user_type"`
	Streak           int                 `json:"streak"`
	HasStreakToday   bool                `json:"has_streak_today"`
	OwnedWorkspaces  []WorkspaceResponse `json:"owned_workspaces"`
	CollabWorkspaces []WorkspaceResponse `json:"collab_workspaces"`
}

// NewProfileResponse converts a User model to a ProfileResponse DTO.
func NewProfileResponse(user *models.User) ProfileResponse {
	owned := make([]WorkspaceResponse, len(user.OwnedWorkspaces))
	for i, w := range user.OwnedWorkspaces {
		owned[i] = NewWorkspaceResponse(&w)
	}

	collab := make([]WorkspaceResponse, 0)
	for _, wUser := range user.CollabWorkspaces {
		if wUser.Workspace.ID != 0 {
			collab = append(collab, NewWorkspaceResponse(&wUser.Workspace))
		}
	}

	return ProfileResponse{
		ID:               user.ID,
		FullName:         user.FullName,
		Email:            user.Email,
		UserType:         user.UserType,
		Streak:           user.Streak,
		HasStreakToday:   utils.HasStreakToday(user),
		OwnedWorkspaces:  owned,
		CollabWorkspaces: collab,
	}
}
