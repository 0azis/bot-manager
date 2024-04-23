package service

import (
	"botmanager/internal/adapter/repo"
	"botmanager/internal/core/telegram"
	"botmanager/internal/core/domain"
	"botmanager/internal/core/port/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type shopService struct {
	store repo.Store
	pool  *telegram.GoroutinesPool
}

func (sc shopService) RunOneBot(c *fiber.Ctx) error {
	var shopCredentials domain.ShopCredentials
	err := c.BodyParser(&shopCredentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}
	// check for already running bot
	if sc.pool.Exists(shopCredentials.ID) {
		return fiber.NewError(400, http.StatusText(400))
	}

	botData, err := sc.store.Shop.Get(shopCredentials.ID)
	if botData.Token == "" {
		return fiber.NewError(404, http.StatusText(404))
	}

	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	goroutine, err := telegram.New(botData, sc.store, sc.pool)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}
	goroutine.Start()

	return fiber.NewError(200, http.StatusText(200))
}

func (sc shopService) StopOneBot(c *fiber.Ctx) error {
	var shopCredentials domain.ShopCredentials
	err := c.BodyParser(&shopCredentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	runGoroutine := sc.pool.Get(shopCredentials.ID)
	if runGoroutine == nil {
		return fiber.NewError(404, http.StatusText(404))
	}

	runGoroutine.Stop()

	return fiber.NewError(200, http.StatusText(200))
}

func NewShopControllers(store repo.Store, pool *telegram.GoroutinesPool) service.ShopService{
	return shopService{
		store: store,
		pool:  pool,
	}
}
