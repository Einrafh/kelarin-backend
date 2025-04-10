package models

import "time"

// Workspace represents a Kanban workspace.
type Workspace struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Title            string    `json:"title"`
	Purpose          string    `json:"purpose"`
	Description      string    `json:"description"`
	WorkspacePicture string    `json:"workspace_picture"` // Stored file path or URL
	WorkspaceBanner  string    `json:"workspace_banner"`  // Stored file path or URL
	OwnerID          uint      `json:"owner_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Owner is the creator of the workspace.
	Owner User `gorm:"foreignKey:OwnerID" json:"owner"`
	// Collaborators contains the list of workspace collaborators.
	Collaborators []WorkspaceUser `json:"collaborators"`
}
