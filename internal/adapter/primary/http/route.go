package http

import (
	"botmanager/internal/adapter/primary/http/controller"
	"botmanager/internal/adapter/secondary/database"
	"botmanager/internal/core/goroutine"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App, store database.Store, pool *goroutine.GoroutinesPool, homeBot string) {
	// init the main API router
	r := app.Group("/v1")

	// plug other groups of routes
	shopRoutes(r, store, pool)
	mailRoutes(r, store, pool)
	messageRoutes(r, store, pool)
	notificationRoutes(r, store, pool, homeBot)
}

func shopRoutes(r fiber.Router, store database.Store, pool *goroutine.GoroutinesPool) {
	controllers := controller.NewShopControllers(pool, store)

	// set routes of bot group
	bot := r.Group("/bot")
	bot.Post("/", controllers.RunOneBot)
	bot.Delete("/", controllers.StopOneBot)
}

func messageRoutes(r fiber.Router, store database.Store, pool *goroutine.GoroutinesPool) {
	controllers := controller.NewMessageControllers(store, pool)

	mail := r.Group("/message")
	mail.Post("/", controllers.SendMessage)
}

func mailRoutes(r fiber.Router, store database.Store, pool *goroutine.GoroutinesPool) {
	controllers := controller.NewMailControllers(store, pool)

	mail := r.Group("/mail")
	mail.Post("/", controllers.SendMail)
}

func notificationRoutes(r fiber.Router, store database.Store, pool *goroutine.GoroutinesPool, homeBot string) {
	controllers := controller.NewNotificationService(store, pool, homeBot)

	notification := r.Group("/notification")
	notification.Post("/", controllers.SendNotification)
}