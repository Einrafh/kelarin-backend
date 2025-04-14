package controllers

import (
	"kelarin-backend/utils"
	"log"
	"strconv"

	"kelarin-backend/dto"
	"kelarin-backend/models"
	"kelarin-backend/repositories"

	"github.com/gofiber/fiber/v2"
)

// CreateAssignee adds a user as an assignee to a card.
func CreateAssignee(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card_id in route"})
	}

	userIDStr := c.FormValue("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
	}

	assignee := models.CardAssignee{
		CardID: uint(cardID),
		UserID: uint(userID),
	}

	if err := repositories.CreateCardAssignee(&assignee); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add assignee"})
	}

	var populatedAssignee models.CardAssignee
	if err := repositories.GetCardAssignee(uint(cardID), uint(userID), &populatedAssignee); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve populated assignee"})
	}

	if err := utils.IncrementStreak(uint(userID)); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	response := dto.NewCardAssigneeResponse(&populatedAssignee)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"assignee": response})
}

// GetAssignees retrieves all assignees for a given card.
func GetAssignees(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card_id"})
	}

	var assignees []models.CardAssignee
	if err := repositories.GetAssigneesByCardID(uint(cardID), &assignees); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch assignees"})
	}

	// Map each assignee to DTO.
	var response []dto.CardAssigneeResponse
	for _, a := range assignees {
		response = append(response, dto.NewCardAssigneeResponse(&a))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"assignees": response})
}

// GetAssignee retrieves a specific assignee by card_id and user_id.
func GetAssignee(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card_id"})
	}
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
	}

	var assignee models.CardAssignee
	if err := repositories.GetCardAssignee(uint(cardID), uint(userID), &assignee); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Assignee not found"})
	}

	response := dto.NewCardAssigneeResponse(&assignee)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"assignee": response})
}

// DeleteAssignee removes an assignee from a card.
func DeleteAssignee(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card_id"})
	}
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
	}

	if err := repositories.DeleteCardAssignee(uint(cardID), uint(userID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove assignee"})
	}

	if err := utils.IncrementStreak(uint(userID)); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Assignee removed successfully"})
}
