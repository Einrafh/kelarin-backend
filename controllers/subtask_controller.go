package controllers

import (
	"kelarin-backend/utils"
	"log"
	"strconv"

	"kelarin-backend/models"
	"kelarin-backend/repositories"

	"github.com/gofiber/fiber/v2"
)

// CreateSubtask creates a new subtask for a given card.
func CreateSubtask(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card ID"})
	}

	title := c.FormValue("title")
	if title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
	}

	subtask := models.Subtask{
		Title:  title,
		CardID: uint(cardID),
	}

	if err := repositories.CreateSubtask(&subtask); err != nil {
		log.Println("Error creating subtask:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create subtask"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"subtask": subtask})
}

// GetSubtasks retrieves all subtasks for a given card.
func GetSubtasks(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card ID"})
	}

	var subtasks []models.Subtask
	if err := repositories.GetSubtasksByCard(uint(cardID), &subtasks); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch subtasks"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"subtasks": subtasks})
}

// GetSubtask retrieves a single subtask by its ID.
func GetSubtask(c *fiber.Ctx) error {
	subtaskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subtask ID"})
	}

	var subtask models.Subtask
	if err := repositories.GetSubtaskByID(uint(subtaskID), &subtask); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Subtask not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"subtask": subtask})
}

// UpdateSubtask updates an existing subtask.
func UpdateSubtask(c *fiber.Ctx) error {
	subtaskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subtask ID"})
	}

	var subtask models.Subtask
	if err := repositories.GetSubtaskByID(uint(subtaskID), &subtask); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Subtask not found"})
	}

	subtask.Title = c.FormValue("title")
	isDoneStr := c.FormValue("is_done")
	if isDoneStr == "true" {
		subtask.IsDone = true
	} else if isDoneStr == "false" {
		subtask.IsDone = false
	}

	if err := repositories.UpdateSubtask(&subtask); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update subtask"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"subtask": subtask})
}

// DeleteSubtask deletes a subtask by its ID.
func DeleteSubtask(c *fiber.Ctx) error {
	subtaskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subtask ID"})
	}

	if err := repositories.DeleteSubtask(uint(subtaskID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete subtask"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Subtask deleted successfully"})
}
