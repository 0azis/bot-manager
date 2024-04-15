package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func InitRoutes(app *fiber.App, store *sqlx.DB) {
	// init the main API router
	r := app.Group("/v1")

	// plug other groups of routes
	shopRoutes(r, store)
}
