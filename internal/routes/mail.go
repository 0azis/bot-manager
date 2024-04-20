package routes

import (
	"botmanager/internal/controllers"
	"botmanager/internal/models/goroutine"
	"botmanager/internal/repos"

	"github.com/gofiber/fiber/v2"
)

func mailRoutes(r fiber.Router, store repos.Store, pool *goroutine.GoroutinesPool) {
	controllers := controllers.NewMailControllers(store, *pool)

	mail := r.Group("/mail")
	mail.Post("/", controllers.Send)
}
