package router

import (
	"goploy/src/apis/handler"

	"github.com/gofiber/fiber/v3"
)

// AuthRouter defines HTTP routes for authentication endpoints.
func AuthRouter(app fiber.Router, handler *handler.AuthHandler) {
	auth := app.Group("/auth")
	auth.Get("/setup", handler.Setup)
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Post("/refresh", handler.Refresh)
	auth.Post("/logout", handler.Logout)
}
