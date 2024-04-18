package routes

import (
	"botmanager/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func InitRoutes(app *fiber.App, store *sqlx.DB, pool *models.GoroutinesPool) {
	// init the main API router
	r := app.Group("/v1")

	// plug other groups of routes
	shopRoutes(r, store, pool)
	mailRoutes(r, store, pool)
}
