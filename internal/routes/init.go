package routes

import (
	"botmanager/internal/models/goroutine"
	"botmanager/internal/repos"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App, store repos.Store, pool *goroutine.GoroutinesPool) {
	// init the main API router
	r := app.Group("/v1")

	// plug other groups of routes
	shopRoutes(r, store, pool)
	mailRoutes(r, store, pool)
}
