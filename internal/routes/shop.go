package routes

import (
	"botmanager/internal/controllers"
	"botmanager/internal/models"
	"botmanager/internal/repos"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func shopRoutes(r fiber.Router, store *sqlx.DB, pool *models.GoroutinesPool) {
	// init the repositories 
	shopRepo := repos.NewShopRepo(store)
	subRepo := repos.NewSubscriberRepo(store)	

	// init the http controllers from using repo
	controllers := controllers.NewShopControllers(shopRepo, subRepo, pool)

	// set routes of bot group
	bot := r.Group("/bot")
	bot.Post("/", controllers.RunOneBot)
	bot.Delete("/", controllers.StopOneBot)
}
