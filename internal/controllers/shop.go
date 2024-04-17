package controllers

import (
	"botmanager/internal/models"
	"botmanager/internal/repos"
	"botmanager/internal/setup"
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
}

// func (sc ShopControllers) RunAfterUpdate(c *fiber.Ctx) error {
// 	var shopCredentials models.ShopCredentials
// 	err := c.BodyParser(&shopCredentials)
// 	if err != nil {
// 		return fiber.NewError(400, http.StatusText(400))
// 	}

// 	if shopCredentials.IsToken() {
// 		return fiber.NewError(404, http.StatusText(404))
// 	}

// 	updatedBot, err := sc.repo.Get(shopCredentials.Token)
// 	if err != nil {
// 		return fiber.NewError(500, http.StatusText(500))
// 	}

// 	return fiber.NewError(200, http.StatusText(200))
// }

func (sc ShopControllers) RunOneBot(c *fiber.Ctx) error {
	var shopCredentials models.ShopCredentials
	err := c.BodyParser(&shopCredentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	// check for already running bot
	if setup.GoroutineExists(shopCredentials.Token) {
		return fiber.NewError(400, http.StatusText(400))
	}

	// check for unvalid token
	if sc.shopRepo.IsTokenValid(shopCredentials.Token) {
		return fiber.NewError(400, http.StatusText(400))
	}

	// create a channel for this goroutine
	ch := make(chan bool)
	fmt.Println(ch)
	// // add goroutine and its channel to map
	setup.Goroutines[shopCredentials.Token] = ch

	// // start goroutine
	go setup.BotWorker(shopCredentials.Token, sc.shopRepo, sc.subRepo)

	ch <- true

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
	ch := setup.Goroutines[shopCredentials.Token]

	// check for empty channel
	if ch == nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	ch <- false

	// delete goroutine from map
	delete(setup.Goroutines, shopCredentials.Token)
	return fiber.NewError(200, http.StatusText(200))
}

func NewShopControllers(shopRepo repos.ShopRepo, subRepo repos.SubscriberRepo) *ShopControllers {
	return &ShopControllers{
		shopRepo: shopRepo,
		subRepo: subRepo,
	}
}
