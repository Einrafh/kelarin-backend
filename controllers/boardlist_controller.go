package controllers

import (
	"kelarin-backend/utils"
	"log"
	"strconv"
	"time"

	"kelarin-backend/models"
	"kelarin-backend/repositories"

	"github.com/gofiber/fiber/v2"
)

// CreateBoardList creates a new board list within a workspace.
func CreateBoardList(c *fiber.Ctx) error {
	workspaceID, err := strconv.Atoi(c.Params("workspace_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid workspace ID"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	role, err := utils.CheckRoleInWorkspace(userID, uint(workspaceID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user role in workspace"})
	}
	if !utils.IsEditorAdminOwner(role) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permission to create board list"})
	}

	title := c.FormValue("title")
	if title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
	}

	list := models.BoardList{
		Title:       title,
		WorkspaceID: uint(workspaceID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := repositories.CreateBoardList(&list); err != nil {
		log.Println("Error creating board list:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create board list"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"list": list})
}

// GetBoardLists returns all board lists for a given workspace.
func GetBoardLists(c *fiber.Ctx) error {
	workspaceID, err := strconv.Atoi(c.Params("workspace_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid workspace ID"})
	}

	var lists []models.BoardList
	if err := repositories.GetBoardListsByWorkspace(uint(workspaceID), &lists); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch board lists"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"lists": lists})
}

// UpdateBoardList updates an existing board list.
func UpdateBoardList(c *fiber.Ctx) error {
	listID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid list ID"})
	}

	var list models.BoardList
	if err := repositories.GetBoardListByID(uint(listID), &list); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "List not found"})
	}

	list.Title = c.FormValue("title")
	list.UpdatedAt = time.Now()

	if err := repositories.UpdateBoardList(&list); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update board list"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"list": list})
}

// DeleteBoardList deletes a board list.
func DeleteBoardList(c *fiber.Ctx) error {
	listID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid list ID"})
	}

	if err := repositories.DeleteBoardList(uint(listID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete board list"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Board list deleted successfully"})
}
