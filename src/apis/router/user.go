package router

import (
	"goploy/src/apis/guard"
	"goploy/src/apis/handler"

	"github.com/gofiber/fiber/v3"
)

// UserRouter defines HTTP routes for user profile endpoints.
func UserRouter(
	app fiber.Router,
	handler *handler.UserHandler,
	guard *guard.Guard,
) {
	user := app.Group("/user", guard.Auth())
	user.Get("/me", handler.Me)
}
