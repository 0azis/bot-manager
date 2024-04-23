package service

import (
	"botmanager/internal/adapter/repo"
	"botmanager/internal/core/telegram"
	"botmanager/internal/core/domain"
	"botmanager/internal/core/port/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type mailService struct {
	store repo.Store
	pool  *telegram.GoroutinesPool
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

	goroutine := mc.pool.Get(shop.ID)
	err = goroutine.SendMail(mail)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	return fiber.NewError(200, http.StatusText(200))
}

func NewMailControllers(store repo.Store, pool *telegram.GoroutinesPool) service.MailService {
	return mailService{
		store: store,
		pool:  pool,
	}
}
