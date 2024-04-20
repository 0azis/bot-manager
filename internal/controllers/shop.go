package controllers

import (
	"botmanager/internal/models"
	"botmanager/internal/models/goroutine"
	"botmanager/internal/repos"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// implement controllers from repo
type ShopControllers struct {
	store repos.Store
	pool  *goroutine.GoroutinesPool
}

func (sc ShopControllers) RunOneBot(c *fiber.Ctx) error {
	var shopCredentials models.ShopCredentials
	err := c.BodyParser(&shopCredentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}
	// check for already running bot
	if sc.pool.Exists(shopCredentials.Token) {
		return fiber.NewError(400, http.StatusText(400))
	}

	// check for unvalid token
	if sc.store.Shop().IsTokenValid(shopCredentials.Token) {
		return fiber.NewError(400, http.StatusText(400))
	}

	botData, err := sc.store.Shop().Get(shopCredentials.Token)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	goroutine, err := goroutine.New(botData, sc.store, sc.pool)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}
	goroutine.Start()

	return fiber.NewError(200, http.StatusText(200))
}

func (sc ShopControllers) StopOneBot(c *fiber.Ctx) error {
	var shopCredentials models.ShopCredentials
	err := c.BodyParser(&shopCredentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	// check for unvalid token
	if sc.store.Shop().IsTokenValid(shopCredentials.Token) {
		return fiber.NewError(400, http.StatusText(400))
	}

	runGoroutine := sc.pool.Get(shopCredentials.Token)
	if runGoroutine == nil {
		return fiber.NewError(404, http.StatusText(404))
	}

	runGoroutine.Stop()

	return fiber.NewError(200, http.StatusText(200))
}

func NewShopControllers(store repos.Store, pool *goroutine.GoroutinesPool) *ShopControllers {
	return &ShopControllers{
		store: store,
		pool:  pool,
	}
}
