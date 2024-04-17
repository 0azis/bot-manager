package repos

import (
	"botmanager/internal/models"

	"github.com/jmoiron/sqlx"
)

type SubscriberRepo interface {
	Insert(sub models.Subscriber) error
	Select() ([]models.Subscriber, error)
}

type subscriber struct {
	db *sqlx.DB
}

func (s subscriber) Insert(sub models.Subscriber) error {
	_, err := s.db.Query("insert into subscriber (username, first_name, last_name, telegram_id, avatar_url, shop_id) values ($1, $2, $3, $4, $5, $6)", sub.UserName, sub.FirstName, sub.LastName, sub.TelegramID, sub.AvatarUrl, sub.ShopID)
	return err
}

func (s subscriber) Select() ([]models.Subscriber, error) {
	var subs []models.Subscriber
	err := s.db.Select(&subs, `select * from subscriber`)
	return subs, err
}

func NewSubscriberRepo(db *sqlx.DB) *subscriber {
	return &subscriber{
		db: db,
	}
}
