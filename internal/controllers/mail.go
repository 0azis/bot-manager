package controllers

import (
	"botmanager/internal/models"
	"botmanager/internal/repos"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type MailControllers struct {
	mailRepo repos.MailRepo
	shopRepo repos.ShopRepo
	pool *models.GoroutinesPool
}

func (mc MailControllers) Send(c *fiber.Ctx) error {
	var credentials models.MailCredentials
	err := c.BodyParser(&credentials)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}
	mail, err := mc.mailRepo.Get(credentials.ID)		
	if mail.ID == "" {
		return fiber.NewError(404, http.StatusText(404))
	}
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	shop, err := mc.shopRepo.GetBy("id", mail.ShopID)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	ch := mc.pool.Get(shop.Token)

	if ch == nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	msg := models.MailType(mail)

	ch <- msg

	return fiber.NewError(200, http.StatusText(200))
}

func NewMailControllers(mailRepo repos.MailRepo, shopRepo repos.ShopRepo, pool *models.GoroutinesPool) *MailControllers {
	return &MailControllers{
		mailRepo: mailRepo,
		shopRepo: shopRepo,
		pool: pool,
	}
}
