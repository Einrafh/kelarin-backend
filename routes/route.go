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
}
