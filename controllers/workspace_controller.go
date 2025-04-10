package controllers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"kelarin-backend/database"
	"kelarin-backend/dto"
	"kelarin-backend/models"
	"kelarin-backend/repositories"

	"github.com/gofiber/fiber/v2"
)

// isImage validates if the file extension is an image.
func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" ||
		ext == ".gif" || ext == ".bmp" || ext == ".webp"
}

// isVideo validates if the file extension is a video.
func isVideo(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".mp4" || ext == ".avi" || ext == ".mov" ||
		ext == ".mkv" || ext == ".webm"
}

// AddWorkspace creates a new workspace and adds collaborators if provided.
// It accepts form-data. If no picture or banner is uploaded, the corresponding fields are set to empty strings.
func AddWorkspace(c *fiber.Ctx) error {
	// Retrieve basic fields from form-data
	title := c.FormValue("title")
	description := c.FormValue("description")
	purpose := c.FormValue("purpose")
	collaborators := c.FormValue("collaborator") // comma separated emails

	if strings.TrimSpace(title) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
	}

	// Get the user ID from context (set by JWT middleware)
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Retrieve owner information
	var owner models.User
	if err := database.DB.First(&owner, userID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Owner user not found"})
	}

	// Set picture and banner paths to empty strings (frontend handles default images)
	picPath := ""
	bannerPath := ""

	// Prepare upload directory
	uploadDir := filepath.Join(".", "uploads")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			log.Println("Error creating uploads directory:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create uploads directory"})
		}
	}

	// Process workspace_picture file if provided
	if filePicture, err := c.FormFile("workspace_picture"); err == nil {
		if !isImage(filePicture.Filename) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace_picture must be an image"})
		}
		// Validate file size based on user type
		if owner.UserType == "regular" && filePicture.Size > 10<<20 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image too large for regular user (max 10MB)"})
		}
		if owner.UserType == "premium" && filePicture.Size > 20<<20 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image too large for premium user (max 20MB)"})
		}
		picPath = filepath.Join(uploadDir, filePicture.Filename)
		if err := c.SaveFile(filePicture, picPath); err != nil {
			log.Println("Error saving workspace picture:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save workspace picture"})
		}
	}

	// Process workspace_banner file if provided
	if fileBanner, err := c.FormFile("workspace_banner"); err == nil {
		// If banner is a video, restrict regular users
		if isVideo(fileBanner.Filename) {
			if owner.UserType == "regular" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Regular users cannot upload videos"})
			}
			if fileBanner.Size > 10<<20 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Video too large for premium user (max 10MB)"})
			}
		} else {
			// Banner is an image
			if !isImage(fileBanner.Filename) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace_banner must be an image or video"})
			}
			if owner.UserType == "regular" && fileBanner.Size > 10<<20 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image too large for regular user (max 10MB)"})
			}
			if owner.UserType == "premium" && fileBanner.Size > 20<<20 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image too large for premium user (max 10MB)"})
			}
		}
		bannerPath = filepath.Join(uploadDir, fileBanner.Filename)
		if err := c.SaveFile(fileBanner, bannerPath); err != nil {
			log.Println("Error saving workspace banner:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save workspace banner"})
		}
	}

	// Create new workspace object
	newWorkspace := models.Workspace{
		Title:            title,
		Description:      description,
		Purpose:          purpose,
		WorkspacePicture: picPath,
		WorkspaceBanner:  bannerPath,
		OwnerID:          userID,
		CreatedAt:        time.Now(),
	}

	// Save workspace to database
	if err := repositories.CreateWorkspace(&newWorkspace); err != nil {
		log.Println("Error creating workspace:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create workspace"})
	}

	// Add the owner as a collaborator with role "owner"
	if err := repositories.AddCollaboratorToWorkspaceWithRole(&newWorkspace, &owner, "owner"); err != nil {
		log.Println("Error adding owner as collaborator:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add owner as collaborator"})
	}

	// If collaborators field is provided (comma-separated emails), add them as "viewer" by default
	if collaborators != "" {
		emailList := strings.Split(collaborators, ",")
		if failedEmails, _ := repositories.AddCollaboratorsByEmails(&newWorkspace, emailList); len(failedEmails) > 0 {
			log.Println("Failed to add some collaborators:", failedEmails)
		}
	}

	// Preload Owner and Collaborators.User for the response
	if err := database.DB.
		Preload("Owner").
		Preload("Collaborators.User").
		First(&newWorkspace, newWorkspace.ID).Error; err != nil {
		log.Println("Error preloading workspace:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load workspace data"})
	}

	response := dto.NewWorkspaceResponse(&newWorkspace)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Workspace created successfully",
		"workspace": response,
	})
}

// GetAllWorkspaces returns all workspaces using the WorkspaceResponse DTO.
func GetAllWorkspaces(c *fiber.Ctx) error {
	var workspaces []models.Workspace
	if err := repositories.GetAllWorkspacesWithOwner(&workspaces); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch workspaces"})
	}

	// Convert each workspace model to DTO response
	var response []dto.WorkspaceResponse
	for _, ws := range workspaces {
		response = append(response, dto.NewWorkspaceResponse(&ws))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"workspaces": response})
}

// GetAccessibleWorkspaces returns the workspaces accessible by the authenticated user,
// using the WorkspaceResponse DTO.
func GetAccessibleWorkspaces(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var workspaces []models.Workspace
	if err := repositories.GetWorkspacesAccessibleByUser(userID, &workspaces); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch workspaces"})
	}

	var response []dto.WorkspaceResponse
	for _, ws := range workspaces {
		response = append(response, dto.NewWorkspaceResponse(&ws))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"workspaces": response})
}

// GetWorkspace returns the workspace details by ID using the WorkspaceResponse DTO.
func GetWorkspace(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	var ws models.Workspace
	if err := repositories.GetWorkspaceByIDWithOwner(workspaceID, &ws); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Workspace not found"})
	}
	response := dto.NewWorkspaceResponse(&ws)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"workspace": response})
}

// UpdateWorkspace updates a workspace's basic fields and collaborator settings.
// It accepts form-data with JSON strings for collaborator changes.
func UpdateWorkspace(c *fiber.Ctx) error {
	// Retrieve workspace by ID with preloaded Owner and Collaborators.User
	id := c.Params("id")
	var ws models.Workspace
	if err := repositories.GetWorkspaceByIDWithOwner(id, &ws); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Workspace not found"})
	}

	// Ensure that the request is made by the owner
	userID, ok := c.Locals("user_id").(uint)
	if !ok || ws.OwnerID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Bind basic fields using DTO
	var updateReq dto.UpdateWorkspaceRequest
	updateReq.Title = c.FormValue("title")
	updateReq.Purpose = c.FormValue("purpose")
	updateReq.Description = c.FormValue("description")

	// Parse collaborator changes from JSON strings in form-data
	if formVal := c.FormValue("add_collaborators"); formVal != "" {
		if err := json.Unmarshal([]byte(formVal), &updateReq.AddCollaborators); err != nil {
			log.Println("Error parsing add_collaborators:", err)
		}
	}
	if formVal := c.FormValue("remove_collaborators"); formVal != "" {
		if err := json.Unmarshal([]byte(formVal), &updateReq.RemoveCollaborators); err != nil {
			log.Println("Error parsing remove_collaborators:", err)
		}
	}
	if formVal := c.FormValue("update_collaborators"); formVal != "" {
		if err := json.Unmarshal([]byte(formVal), &updateReq.UpdateCollaborators); err != nil {
			log.Println("Error parsing update_collaborators:", err)
		}
	}

	// Update basic fields if provided
	if updateReq.Title != "" {
		ws.Title = updateReq.Title
	}
	if updateReq.Purpose != "" {
		ws.Purpose = updateReq.Purpose
	}
	if updateReq.Description != "" {
		ws.Description = updateReq.Description
	}

	// Prepare upload directory
	uploadDir := filepath.Join(".", "uploads")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return err
	}

	// Update workspace_picture if file is provided
	if filePicture, err := c.FormFile("workspace_picture"); err == nil {
		if !isImage(filePicture.Filename) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace_picture must be an image"})
		}
		picPath := filepath.Join(uploadDir, filePicture.Filename)
		if err := c.SaveFile(filePicture, picPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save workspace_picture"})
		}
		ws.WorkspacePicture = picPath
	}

	// Update workspace_banner if file is provided
	if fileBanner, err := c.FormFile("workspace_banner"); err == nil {
		if !isImage(fileBanner.Filename) && !isVideo(fileBanner.Filename) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace_banner must be an image or video"})
		}
		bannerPath := filepath.Join(uploadDir, fileBanner.Filename)
		if err := c.SaveFile(fileBanner, bannerPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save workspace_banner"})
		}
		ws.WorkspaceBanner = bannerPath
	}

	// Update the updated time and save basic changes
	ws.UpdatedAt = time.Now()
	if err := database.DB.Save(&ws).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update workspace"})
	}

	// ---- Handling Collaborator Changes ----

	// 1. Add new collaborators
	for _, collab := range updateReq.AddCollaborators {
		user, err := repositories.GetUserByEmail(collab.Email)
		if err != nil || user.ID == ws.OwnerID {
			continue
		}
		// Skip if already a collaborator
		if exists, _ := repositories.IsUserAlreadyCollaborator(ws.ID, user.ID); exists {
			continue
		}
		if err := repositories.AddCollaboratorToWorkspaceWithRole(&ws, user, collab.Role); err != nil {
			log.Println("Failed to add collaborator:", collab.Email, err)
		}
	}

	// 2. Remove collaborators
	for _, email := range updateReq.RemoveCollaborators {
		email = strings.TrimSpace(email)
		user, err := repositories.GetUserByEmail(email)
		if err != nil || user.ID == ws.OwnerID {
			continue
		}
		if err := database.DB.
			Where("workspace_id = ? AND user_id = ?", ws.ID, user.ID).
			Delete(&models.WorkspaceUser{}).Error; err != nil {
			log.Println("Failed to remove collaborator:", email, err)
		}
	}

	// 3. Update collaborator roles
	for _, collab := range updateReq.UpdateCollaborators {
		user, err := repositories.GetUserByEmail(collab.Email)
		if err != nil {
			continue
		}
		if err := database.DB.
			Model(&models.WorkspaceUser{}).
			Where("workspace_id = ? AND user_id = ?", ws.ID, user.ID).
			Update("role", collab.Role).Error; err != nil {
			log.Println("Failed to update collaborator role for:", collab.Email, err)
		}
	}

	// Reload workspace with preloaded data for response
	if err := repositories.GetWorkspaceByIDWithOwner(id, &ws); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load updated workspace"})
	}

	response := dto.NewWorkspaceResponse(&ws)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Workspace updated successfully",
		"workspace": response,
	})
}

// ShareWorkspace shares a workspace with a user by adding them as a collaborator with a specified role.
// Only the workspace owner is allowed to share the workspace.
func ShareWorkspace(c *fiber.Ctx) error {
	// Retrieve workspace ID from URL parameters
	workspaceID := c.Params("id")

	// Retrieve form-data values for email and role
	email := c.FormValue("email")
	role := c.FormValue("role")

	// Validate that email and role are provided
	if email == "" || role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both email and role are required in the form-data.",
		})
	}

	// Create the payload DTO from form-data
	payload := dto.ShareWorkspaceRequest{
		Email: email,
		Role:  role,
	}

	// Retrieve the workspace details, including its owner
	var ws models.Workspace
	if err := repositories.GetWorkspaceByIDWithOwner(workspaceID, &ws); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Workspace not found"})
	}

	// Only allow the owner of the workspace to share it
	userID := c.Locals("user_id").(uint)
	if ws.OwnerID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not the owner of this workspace"})
	}

	// Find the user by the provided email
	user, err := repositories.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if the user is already a collaborator in the workspace
	if exists, _ := repositories.IsUserAlreadyCollaborator(ws.ID, user.ID); exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already a collaborator"})
	}

	// Add the user as a collaborator with the specified role from the payload
	if err := repositories.AddCollaboratorToWorkspaceWithRole(&ws, user, payload.Role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add collaborator"})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Workspace shared successfully"})
}

// DeleteWorkspace deletes a workspace by its ID.
func DeleteWorkspace(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	var ws models.Workspace
	if err := repositories.GetWorkspaceByIDWithOwner(workspaceID, &ws); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Workspace not found"})
	}

	if err := repositories.DeleteWorkspace(workspaceID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete workspace"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Workspace deleted successfully"})
}
