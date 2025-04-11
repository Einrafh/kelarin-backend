package controllers

import (
	"log"
	"strconv"
	"time"

	"kelarin-backend/models"
	"kelarin-backend/repositories"

	"github.com/gofiber/fiber/v2"
)

// CreateCard creates a new card in a specific board list.
func CreateCard(c *fiber.Ctx) error {
	listID, err := strconv.Atoi(c.Params("list_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid list ID"})
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	deadlineStr := c.FormValue("deadline") // optional deadline field

	var deadline *time.Time
	if deadlineStr != "" {
		parsed, err := time.Parse(time.RFC3339, deadlineStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid deadline format"})
		}
		deadline = &parsed
	}

	card := models.Card{
		Title:       title,
		Description: description,
		Deadline:    deadline,
		ListID:      uint(listID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := repositories.CreateCard(&card); err != nil {
		log.Println("Error creating card:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create card"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"card": card})
}

// GetCards retrieves all cards for a given list.
func GetCards(c *fiber.Ctx) error {
	listID, err := strconv.Atoi(c.Params("list_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid list ID"})
	}

	var cards []models.Card
	if err := repositories.GetCardsByListID(uint(listID), &cards); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch cards"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"cards": cards})
}

// GetCard returns a card by its ID including attachments, labels, and comments.
func GetCard(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card ID"})
	}

	var card models.Card
	if err := repositories.GetCardByID(uint(cardID), &card); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Card not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"card": card})
}

// UpdateCard updates a card's basic fields.
func UpdateCard(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card ID"})
	}

	var card models.Card
	if err := repositories.GetCardByID(uint(cardID), &card); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Card not found"})
	}

	card.Title = c.FormValue("title")
	card.Description = c.FormValue("description")
	deadlineStr := c.FormValue("deadline")
	if deadlineStr != "" {
		parsed, err := time.Parse(time.RFC3339, deadlineStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid deadline format"})
		}
		card.Deadline = &parsed
	}
	card.UpdatedAt = time.Now()

	if err := repositories.UpdateCard(&card); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update card"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"card": card})
}

// DeleteCard deletes a card by its ID.
func DeleteCard(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card ID"})
	}

	if err := repositories.DeleteCard(uint(cardID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete card"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Card deleted successfully"})
}
