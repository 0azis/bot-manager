package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func InitRoutes(app *fiber.App, store *sqlx.DB) {
	// init the main API router
	_ = app.Group("/v1")	
}
