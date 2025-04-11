package dto

import "kelarin-backend/models"

// CardAssigneeResponse represents an assignee on a card with user info mapped using ProfileResponse.
type CardAssigneeResponse struct {
	CardID uint            `json:"card_id"`
	UserID uint            `json:"user_id"`
	User   ProfileResponse `json:"user"`
}

// NewCardAssigneeResponse converts a CardAssignee model into a CardAssigneeResponse.
func NewCardAssigneeResponse(assignee *models.CardAssignee) CardAssigneeResponse {
	return CardAssigneeResponse{
		CardID: assignee.CardID,
		UserID: assignee.UserID,
		User:   NewProfileResponse(&assignee.User),
	}
}
