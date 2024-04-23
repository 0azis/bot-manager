package service

import (
	"botmanager/internal/adapter/repo"
	"botmanager/internal/core/telegram"
	"botmanager/internal/core/domain"
	"botmanager/internal/core/port/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type messageService struct {
	store repo.Store
	pool  *telegram.GoroutinesPool
}

func (mc messageService) SendMessage(c *fiber.Ctx) error {
	var messageCredentials domain.MessageCredentials
	err := c.BodyParser(&messageCredentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	message, err := mc.store.Message.Get(messageCredentials.ID)
	if message.ID == "" {
		return fiber.NewError(404, http.StatusText(404))
	}
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	shop, err := mc.store.Shop.Get(message.BotID)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	goroutine := mc.pool.Get(shop.ID)

	err = goroutine.SendMessage(message)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	return fiber.NewError(200, http.StatusText(200))
}

func NewMessageControllers(store repo.Store, pool *telegram.GoroutinesPool) service.MessageService {
	return messageService{
		store: store,
		pool:  pool,
	}
}
