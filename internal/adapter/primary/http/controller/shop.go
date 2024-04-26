package controller

import (
	"botmanager/internal/adapter/secondary/database"
	"botmanager/internal/adapter/secondary/redis"
	"botmanager/internal/core/domain"
	"botmanager/internal/core/goroutine"
	"botmanager/internal/core/port/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type shopService struct {
	pool    *goroutine.GoroutinesPool
	store   database.Store
	redisDB redis.RedisInterface
}

func (sc shopService) RunOneBot(c *fiber.Ctx) error {
	var shopCredentials domain.ShopCredentials
	err := c.BodyParser(&shopCredentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}
	// check for already running bot
	if sc.pool.Exists(shopCredentials.Token) {
		return fiber.NewError(400, http.StatusText(400))
	}

	botData, err := sc.store.Shop.Get(shopCredentials.Token)
	if botData.Token == "" {
		return fiber.NewError(404, http.StatusText(404))
	}

	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	shopBot, err := goroutine.New(botData.Token, sc.pool, sc.store, sc.redisDB)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}
	shopBot.InitShopHandlers()
	shopBot.Start()

	return fiber.NewError(200, http.StatusText(200))
}

func (sc shopService) StopOneBot(c *fiber.Ctx) error {
	var shopCredentials domain.ShopCredentials
	err := c.BodyParser(&shopCredentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	shopBot := sc.pool.Get(shopCredentials.Token)
	if shopBot == nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	shopBot.Stop()
	return fiber.NewError(200, http.StatusText(200))
}

func NewShopControllers(pool *goroutine.GoroutinesPool, store database.Store) service.ShopService {
	return shopService{
		pool:  pool,
		store: store,
	}
}
