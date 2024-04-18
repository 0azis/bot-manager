package controllers

import (
	"botmanager/internal/models"
	"botmanager/internal/repos"
	"botmanager/internal/tools"
	"fmt"
	"net/http"

	// "time"

	// "time"

	"github.com/gofiber/fiber/v2"
)

// implement controllers from repo
type ShopControllers struct {
	shopRepo repos.ShopRepo
	subRepo  repos.SubscriberRepo
	pool *models.GoroutinesPool
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
	if sc.shopRepo.IsTokenValid(shopCredentials.Token) {
		return fiber.NewError(400, http.StatusText(400))
	}

	// create a channel for this goroutine
	ch := make(chan models.ChannelMessage)
	fmt.Println(ch)
	// // add goroutine and its channel to map
	sc.pool.Add(shopCredentials.Token, ch)

	// // start goroutine
	go tools.BotWorker(shopCredentials.Token, sc.shopRepo, sc.subRepo, sc.pool)

	msg := models.WorkType(true)
	ch <- msg 

	return fiber.NewError(200, http.StatusText(200))
}

func (sc ShopControllers) StopOneBot(c *fiber.Ctx) error {
	var shopCredentials models.ShopCredentials
	err := c.BodyParser(&shopCredentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	// check for unvalid token
	if sc.shopRepo.IsTokenValid(shopCredentials.Token) {
		return fiber.NewError(400, http.StatusText(400))
	}

	// get channel from map
	ch := sc.pool.Get(shopCredentials.Token)

	// check for empty channel
	if ch == nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	msg := models.WorkType(false) 

	ch <- msg 

	// delete goroutine from map
	sc.pool.Delete(shopCredentials.Token)
	return fiber.NewError(200, http.StatusText(200))
}

func NewShopControllers(shopRepo repos.ShopRepo, subRepo repos.SubscriberRepo, pool *models.GoroutinesPool) *ShopControllers {
	return &ShopControllers{
		shopRepo: shopRepo,
		subRepo:  subRepo,
		pool: pool,
	}
}
