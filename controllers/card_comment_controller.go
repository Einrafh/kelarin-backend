package controllers

import (
	"kelarin-backend/utils"
	"log"
	"strconv"
	"time"

	"kelarin-backend/dto"
	"kelarin-backend/models"
	"kelarin-backend/repositories"

	"github.com/gofiber/fiber/v2"
)

// CreateCardComment adds a comment to a card.
// The user_id is automatically retrieved from the context.
func CreateCardComment(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card ID in route"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	commentText := c.FormValue("comment")
	if commentText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Comment cannot be empty"})
	}

	comment := models.CardComment{
		CardID:    uint(cardID),
		UserID:    userID,
		Comment:   commentText,
		CreatedAt: time.Now(),
	}

	if err := repositories.CreateCardComment(&comment); err != nil {
		log.Println("Error creating card comment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create card comment"})
	}

	var populatedComment models.CardComment
	if err := repositories.GetCardCommentByID(comment.ID, &populatedComment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to preload comment data"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	response := dto.NewCardCommentResponse(&populatedComment)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"comment": response})
}

// GetComments retrieves all comments for a given card.
func GetComments(c *fiber.Ctx) error {
	cardID, err := strconv.Atoi(c.Params("card_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid card_id"})
	}

	var comments []models.CardComment
	if err := repositories.GetCommentsByCardID(uint(cardID), &comments); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch comments"})
	}

	var response []dto.CardCommentResponse
	for _, comment := range comments {
		response = append(response, dto.NewCardCommentResponse(&comment))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"comments": response})
}

// GetCardComment retrieves a card comment by its ID.
func GetCardComment(c *fiber.Ctx) error {
	commentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}

	var comment models.CardComment
	if err := repositories.GetCardCommentByID(uint(commentID), &comment); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
	}

	response := dto.NewCardCommentResponse(&comment)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"comment": response})
}

// UpdateCardComment updates a card comment.
func UpdateCardComment(c *fiber.Ctx) error {
	commentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}

	var comment models.CardComment
	if err := repositories.GetCardCommentByID(uint(commentID), &comment); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
	}

	newComment := c.FormValue("comment")
	if newComment != "" {
		comment.Comment = newComment
		comment.CreatedAt = time.Now() // Optional: update timestamp
	}

	if err := repositories.UpdateCardComment(&comment); err != nil {
		log.Println("Error updating card comment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update comment"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	response := dto.NewCardCommentResponse(&comment)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"comment": response})
}

// DeleteCardComment deletes a card comment by its ID.
func DeleteCardComment(c *fiber.Ctx) error {
	commentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}
	if err := repositories.DeleteCardComment(uint(commentID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete comment"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := utils.IncrementStreak(userID); err != nil {
		log.Println("Error incrementing streak:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Comment deleted successfully"})
}
