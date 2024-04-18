package routes

import (
	"botmanager/internal/controllers"
	"botmanager/internal/models"
	"botmanager/internal/repos"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func mailRoutes(r fiber.Router, store *sqlx.DB, pool *models.GoroutinesPool) {
	mailRepo := repos.NewMailRepo(store)
	shopRepo := repos.NewShopRepo(store)

	controllers := controllers.NewMailControllers(mailRepo, shopRepo, pool)

	mail := r.Group("/mail")
	mail.Post("/", controllers.Send)
}
