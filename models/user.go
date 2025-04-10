package models

// User represents an application user.
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	FullName string `gorm:"not null;size:255" json:"fullname"`
	Email    string `gorm:"unique;not null;size:255" json:"email"`
	Password string `gorm:"not null;size:255" json:"password"`
	// UserType can be "regular" or "premium"
	UserType string `gorm:"not null;default:'regular'" json:"user_type"`

	// OwnedWorkspaces are the workspaces that the user owns.
	OwnedWorkspaces []Workspace `gorm:"foreignKey:OwnerID" json:"owned_workspaces"`
	// CollabWorkspaces are the workspaces that the user collaborates on.
	CollabWorkspaces []WorkspaceUser `gorm:"foreignKey:UserID" json:"collab_workspaces"`
}
