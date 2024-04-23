package route

import (
	"botmanager/internal/adapter/http/service"
	"botmanager/internal/adapter/repo"
	"botmanager/internal/core/telegram"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App, store repo.Store, pool *telegram.GoroutinesPool) {
	// init the main API router
	r := app.Group("/v1")

	// plug other groups of routes
	shopRoutes(r, store, pool)
	mailRoutes(r, store, pool)
	messageRoutes(r, store, pool)
}

func shopRoutes(r fiber.Router, store repo.Store, pool *telegram.GoroutinesPool) {
	controllers := service.NewShopControllers(store, pool)

	// set routes of bot group
	bot := r.Group("/bot")
	bot.Post("/", controllers.RunOneBot)
	bot.Delete("/", controllers.StopOneBot)
}

func messageRoutes(r fiber.Router, store repo.Store, pool *telegram.GoroutinesPool) {
	controllers := service.NewMessageControllers(store, pool)

	mail := r.Group("/message")
	mail.Post("/", controllers.SendMessage)
}

func mailRoutes(r fiber.Router, store repo.Store, pool *telegram.GoroutinesPool) {
	controllers := service.NewMailControllers(store, pool)

	mail := r.Group("/mail")
	mail.Post("/", controllers.SendMail)
}
