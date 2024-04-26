package repository

import "botmanager/internal/core/domain"

type ShopRepository interface {
	Select() ([]domain.Shop, error)
	Get(ID string) (domain.Shop, error)
	GetByBotID(botID string) (domain.Shop, error)
}

type SubscriberRepository interface {
	Insert(sub domain.Subscriber) error
	Select(shopID string) ([]domain.Subscriber, error)
	IsSubscribed(tgID, shopID string) bool
	Get(ID string) (domain.Subscriber, error)
	GetByTelegramID(telegramID string) (domain.Subscriber, error)
}

type MessageRepository interface {
	Insert(msg domain.Message) error
	Get(ID string) (domain.Message, error)
}

type MailRepository interface {
	Get(mailID string) (domain.Mail, error)
}

type UserRepository interface {
	Insert(user domain.User) error
}
