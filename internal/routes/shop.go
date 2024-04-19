package routes

import (
	"botmanager/internal/controllers"
	"botmanager/internal/models/goroutine"
	"botmanager/internal/repos"

	"github.com/gofiber/fiber/v2"
)

func shopRoutes(r fiber.Router, store repos.Store, pool goroutine.GoroutinesPool) {
	controllers := controllers.NewShopControllers(store, pool)

	// set routes of bot group
	bot := r.Group("/bot")
	bot.Post("/", controllers.RunOneBot)
	bot.Delete("/", controllers.StopOneBot)
}
