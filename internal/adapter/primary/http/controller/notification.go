package controller

import (
	"botmanager/internal/adapter/secondary/database"
	"botmanager/internal/core/domain"
	"botmanager/internal/core/goroutine"
	"botmanager/internal/core/port/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type notificationService struct {
	store   database.Store
	pool    *goroutine.GoroutinesPool
	homeBot string
}

func (ns notificationService) SendNotification(c *fiber.Ctx) error {
	var notification domain.Notification
	err := c.BodyParser(&notification)
	if err != nil {
		return fiber.NewError(400, http.StatusText(400))
	}

	goroutine := ns.pool.Get(ns.homeBot)
	err = goroutine.SendNotification(notification)
	if err != nil {
		return fiber.NewError(500, http.StatusText(500))
	}

	return fiber.NewError(200, http.StatusText(200))
}

func NewNotificationService(store database.Store, pool *goroutine.GoroutinesPool, homeBot string) service.NotificationService {
	return &notificationService{
		store:   store,
		pool:    pool,
		homeBot: homeBot,
	}
}
