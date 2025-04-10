package dto

// UpdateCollaborator represents a collaborator update payload.
type UpdateCollaborator struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// UpdateWorkspaceRequest represents the payload for updating a workspace.
// It is used to bind form-data (with JSON strings for collaborator fields).
type UpdateWorkspaceRequest struct {
	Title               string               `json:"title"`
	Purpose             string               `json:"purpose"`
	Description         string               `json:"description"`
	AddCollaborators    []UpdateCollaborator `json:"add_collaborators"`
	RemoveCollaborators []string             `json:"remove_collaborators"`
	UpdateCollaborators []UpdateCollaborator `json:"update_collaborators"`
}
