package models

// WorkspaceUser is the join model between Workspace and User.
// It defines the role of a user in a workspace (viewer, editor, admin, owner).
type WorkspaceUser struct {
	UserID      uint   `gorm:"primaryKey" json:"user_id"`
	WorkspaceID uint   `gorm:"primaryKey" json:"workspace_id"`
	Role        string `gorm:"not null;default:'viewer'" json:"role"`

	// User contains the user details.
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	// Workspace contains the workspace details.
	Workspace Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"workspace"`
}
