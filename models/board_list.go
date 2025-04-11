package models

import "time"

// BoardList represents a list/column on a Kanban board (e.g., To Do, In Progress).
type BoardList struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	WorkspaceID uint      `gorm:"not null" json:"workspace_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// The workspace this list belongs to.
	Workspace Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"-"`
	// Cards in this list.
	Cards []Card `gorm:"foreignKey:ListID" json:"cards"`
}
