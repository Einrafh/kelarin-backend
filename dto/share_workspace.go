package dto

// ShareWorkspaceRequest represents the payload to share a workspace with a specific role.
// Email is the user's email to be added as a collaborator.
// Role is the role that the user will have in the workspace (e.g., viewer, editor, admin).
type ShareWorkspaceRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required"`
}
