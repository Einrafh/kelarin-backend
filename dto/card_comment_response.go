package dto

import (
	"time"

	"kelarin-backend/models"
)

// CardCommentResponse represents a comment on a card with user info mapped using ProfileResponse.
type CardCommentResponse struct {
	ID        uint            `json:"id"`
	CardID    uint            `json:"card_id"`
	User      ProfileResponse `json:"user"`
	Comment   string          `json:"comment"`
	CreatedAt time.Time       `json:"created_at"`
}

// NewCardCommentResponse converts a CardComment model into a CardCommentResponse.
func NewCardCommentResponse(comment *models.CardComment) CardCommentResponse {
	return CardCommentResponse{
		ID:        comment.ID,
		CardID:    comment.CardID,
		User:      NewProfileResponse(&comment.User),
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
	}
}
