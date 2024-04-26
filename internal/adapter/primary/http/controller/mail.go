package controller

import (
	"botmanager/internal/adapter/secondary/database"
	"botmanager/internal/core/domain"
	"botmanager/internal/core/goroutine"
	"botmanager/internal/core/port/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type mailService struct {
	store database.Store
	pool  *goroutine.GoroutinesPool
}

func (mc mailService) SendMail(c *fiber.Ctx) error {
	var credentials domain.MailCredentials
	err := c.BodyParser(&credentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	mail, err := mc.store.Mail.Get(credentials.ID)
	if mail.ID == "" {
		return fiber.NewError(404, http.StatusText(404))
	}
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	shop, err := mc.store.Shop.Get(mail.ShopID)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	goroutine := mc.pool.Get(shop.Token)
	err = goroutine.SendMail(mail)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	return fiber.NewError(200, http.StatusText(200))
}

func NewMailControllers(store database.Store, pool *goroutine.GoroutinesPool) service.MailService {
	return mailService{
		store: store,
		pool:  pool,
	}
}
