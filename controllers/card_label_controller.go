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

// CreateCardLabel adds a label to a card.
func CreateCardLabel(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card ID in route"})
	}

	name := c.FormValue("name")
	color := c.FormValue("color") // optional

	label := models.CardLabel{
		CardID:    uint(cardID),
		Name:      name,
		Color:     color,
		CreatedAt: time.Now(),
	}

	if err := repositories.CreateCardLabel(&label); err != nil {
		log.Println("Error creating card label:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create card label"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"label": label})
}

// GetLabels retrieves all labels for a given card.
func GetLabels(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card_id"})
	}

	var labels []models.CardLabel
	if err := repositories.GetLabelsByCardID(uint(cardID), &labels); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch labels"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"labels": labels})
}

// GetCardLabel retrieves a card label by its ID.
func GetCardLabel(c *fiber.Ctx) error {
	labelID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid label ID"})
	}

	var label models.CardLabel
	if err := repositories.GetCardLabelByID(uint(labelID), &label); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Label not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"label": label})
}

// UpdateCardLabel updates a card label.
func UpdateCardLabel(c *fiber.Ctx) error {
	labelID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid label ID"})
	}

	var label models.CardLabel
	if err := repositories.GetCardLabelByID(uint(labelID), &label); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Label not found"})
	}

	newName := c.FormValue("name")
	newColor := c.FormValue("color")
	if newName != "" {
		label.Name = newName
	}
	if newColor != "" {
		label.Color = newColor
	}

	if err := repositories.UpdateCardLabel(&label); err != nil {
		log.Println("Error updating card label:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update label"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"label": label})
}

// DeleteCardLabel deletes a card label by its ID.
func DeleteCardLabel(c *fiber.Ctx) error {
	labelID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid label ID"})
	}
	if err := repositories.DeleteCardLabel(uint(labelID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete label"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Label deleted successfully"})
}
