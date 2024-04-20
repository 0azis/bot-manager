package controllers

import (
	"botmanager/internal/models"
	"botmanager/internal/models/goroutine"
	"botmanager/internal/repos"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type MailControllers struct {
	store repos.Store
	pool  goroutine.GoroutinesPool
}

func (mc MailControllers) Send(c *fiber.Ctx) error {
	var credentials models.MailCredentials
	err := c.BodyParser(&credentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	mail, err := mc.store.Mail().Get(credentials.ID)
	if mail.ID == "" {
		return fiber.NewError(404, http.StatusText(404))
	}
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	shop, err := mc.store.Shop().GetBy("id", mail.ShopID)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	goroutine := mc.pool.Get(shop.Token)
	err = goroutine.SendMessages(mail)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	return fiber.NewError(200, http.StatusText(200))
}

func NewMailControllers(store repos.Store, pool goroutine.GoroutinesPool) *MailControllers {
	return &MailControllers{
		store: store,
		pool:  pool,
	}
}
