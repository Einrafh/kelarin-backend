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

// CreateCardAttachment adds an attachment to a card.
// Expects form-data: "card_id", "url" (or file upload handling can be added), "file_name" optional.
func CreateCardAttachment(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card ID in route"})
	}

	url := c.FormValue("url")
	fileName := c.FormValue("file_name") // optional

	attachment := models.CardAttachment{
		CardID:    uint(cardID),
		URL:       url,
		FileName:  fileName,
		CreatedAt: time.Now(),
	}

	if err := repositories.CreateCardAttachment(&attachment); err != nil {
		log.Println("Error creating card attachment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create card attachment"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"attachment": attachment})
}

// GetAttachments retrieves all attachments for a given card.
func GetAttachments(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card_id"})
	}

	var attachments []models.CardAttachment
	if err := repositories.GetAttachmentsByCardID(uint(cardID), &attachments); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch attachments"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"attachments": attachments})
}

// GetCardAttachment retrieves a card attachment by its ID.
func GetCardAttachment(c *fiber.Ctx) error {
	attachmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid attachment ID"})
	}

	var attachment models.CardAttachment
	if err := repositories.GetCardAttachmentByID(uint(attachmentID), &attachment); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Attachment not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"attachment": attachment})
}

// UpdateCardAttachment updates an existing card attachment.
func UpdateCardAttachment(c *fiber.Ctx) error {
	attachmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid attachment ID"})
	}

	var attachment models.CardAttachment
	if err := repositories.GetCardAttachmentByID(uint(attachmentID), &attachment); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Attachment not found"})
	}

	newURL := c.FormValue("url")
	newFileName := c.FormValue("file_name")
	if newURL != "" {
		attachment.URL = newURL
	}
	if newFileName != "" {
		attachment.FileName = newFileName
	}

	if err := repositories.UpdateCardAttachment(&attachment); err != nil {
		log.Println("Error updating card attachment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update attachment"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"attachment": attachment})
}

// DeleteCardAttachment deletes a card attachment by its ID.
func DeleteCardAttachment(c *fiber.Ctx) error {
	attachmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid attachment ID"})
	}

	if err := repositories.DeleteCardAttachment(uint(attachmentID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete attachment"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Attachment deleted successfully"})
}
