package router

import (
	"goploy/src/apis/handler"

	"github.com/gofiber/fiber/v3"
)

// HealthRouter defines HTTP routes for system status endpoints.
func HealthRouter(app fiber.Router, handler *handler.HealthHandler) {
	app.Get("/ping", handler.Ping)
	app.Get("/pong", handler.Pong)
	app.Get("/health", handler.Health)
}
