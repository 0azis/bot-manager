package service

import "github.com/gofiber/fiber/v2"

type ShopService interface {
	RunOneBot(c *fiber.Ctx) error
	StopOneBot(c *fiber.Ctx) error
}

type MessageService interface {
	SendMessage(c *fiber.Ctx) error
}

type MailService interface {
	SendMail(c *fiber.Ctx) error
}

type NotificationService interface {
	SendNotification(c *fiber.Ctx) error
}