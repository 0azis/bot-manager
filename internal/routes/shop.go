package routes

import (
	"botmanager/internal/controllers"
	"botmanager/internal/repos"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func shopRoutes(r fiber.Router, store *sqlx.DB) {
	// init the repository of the shop
	repo := repos.NewShopRepo(store)

	// init the http controllers from using repo
	controllers := controllers.NewShopControllers(repo)

	// set routes of bot group
	bot := r.Group("/bot")
	bot.Post("/", controllers.RunOneBot)
	// bot.Patch("/", controllers.RunAfterUpdate)
	bot.Delete("/", controllers.StopOneBot)
}
