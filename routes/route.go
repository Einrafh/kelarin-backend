package routes

import (
	"github.com/gofiber/fiber/v2"
	"kelarin-backend/controllers"
	"kelarin-backend/middleware"
)

// SetupRoutes initializes all API routes.
func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Welcome to KelarIn Backend API!")
	})

	api := app.Group("/api")

	// Auth routes
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/profile", middleware.AuthMiddleware, controllers.GetProfile)

	// Workspace routes
	workspace := api.Group("/workspace", middleware.AuthMiddleware)
	workspace.Post("/", controllers.AddWorkspace)                     // Create workspace
	workspace.Post("/:id/share", controllers.ShareWorkspace)          // Share workspace
	workspace.Get("/all", controllers.GetAllWorkspaces)               // Get all workspaces
	workspace.Get("/accessible", controllers.GetAccessibleWorkspaces) // Get accessible workspaces
	workspace.Get("/:id", controllers.GetWorkspace)                   // Get workspace by ID
	workspace.Put("/:id", controllers.UpdateWorkspace)                // Update workspace
	workspace.Delete("/:id", controllers.DeleteWorkspace)             // Delete workspace

	// Kanban Board routes
	kanban := api.Group("/kanban", middleware.AuthMiddleware)

	// BoardList routes:
	kanban.Post("/workspace/:workspace_id/lists", controllers.CreateBoardList)
	kanban.Get("/workspace/:workspace_id/lists", controllers.GetBoardLists)
	kanban.Put("/lists/:id", controllers.UpdateBoardList)
	kanban.Delete("/lists/:id", controllers.DeleteBoardList)

	// Card routes:
	kanban.Post("/lists/:list_id/cards", controllers.CreateCard)
	kanban.Get("/lists/:list_id/cards", controllers.GetCards)
	kanban.Get("/cards/:id", controllers.GetCard)
	kanban.Put("/cards/:id", controllers.UpdateCard)
	kanban.Delete("/cards/:id", controllers.DeleteCard)

	// Card Assignee routes:
	kanban.Post("/cards/:card_id/assignees", controllers.CreateAssignee)
	kanban.Get("/cards/:card_id/assignees", controllers.GetAssignees)
	kanban.Get("/cards/:card_id/assignees/:user_id", controllers.GetAssignee)
	kanban.Delete("/cards/:card_id/assignees/:user_id", controllers.DeleteAssignee)

	// Card Label routes:
	kanban.Post("/cards/:card_id/label", controllers.CreateCardLabel)
	kanban.Get("/cards/:card_id/labels", controllers.GetLabels)
	kanban.Get("/cards/label/:id", controllers.GetCardLabel)
	kanban.Put("/cards/label/:id", controllers.UpdateCardLabel)
	kanban.Delete("/cards/label/:id", controllers.DeleteCardLabel)

	// Card Attachment routes:
	kanban.Post("/cards/:card_id/attachment", controllers.CreateCardAttachment)
	kanban.Get("/cards/:card_id/attachments", controllers.GetAttachments)
	kanban.Get("/cards/attachment/:id", controllers.GetCardAttachment)
	kanban.Put("/cards/attachment/:id", controllers.UpdateCardAttachment)
	kanban.Delete("/cards/attachment/:id", controllers.DeleteCardAttachment)

	// Card Comment routes:
	kanban.Post("/cards/:card_id/comment", controllers.CreateCardComment)
	kanban.Get("/cards/:card_id/comments", controllers.GetComments)
	kanban.Get("/cards/comment/:id", controllers.GetCardComment)
	kanban.Put("/cards/comment/:id", controllers.UpdateCardComment)
	kanban.Delete("/cards/comment/:id", controllers.DeleteCardComment)

	// Subtask routes:
	kanban.Post("/cards/:card_id/subtask", controllers.CreateSubtask)
	kanban.Get("/cards/:card_id/subtasks", controllers.GetSubtasks)
	kanban.Get("/subtask/:id", controllers.GetSubtask)
	kanban.Put("/subtask/:id", controllers.UpdateSubtask)
	kanban.Delete("/subtask/:id", controllers.DeleteSubtask)
}
