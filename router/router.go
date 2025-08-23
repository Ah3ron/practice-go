package router

import (
	"app/handler"
	"app/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)

	// User
	user := api.Group("/user")
	user.Get("/", handler.GetAllUsers)
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// Resource
	resource := api.Group("/resource")
	resource.Get("/", handler.GetAllResources)
	resource.Get("/:id", handler.GetResource)
	resource.Post("/", middleware.Protected(), handler.CreateResource)
	resource.Put("/:id", middleware.Protected(), handler.UpdateResource)
	resource.Delete("/:id", middleware.Protected(), handler.DeleteResource)
	resource.Get("/:id/history", handler.GetResourceHistory)

}
