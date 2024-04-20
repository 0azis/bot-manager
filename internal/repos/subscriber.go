package repos

import (
	"botmanager/internal/models"

	"github.com/jmoiron/sqlx"
)

type subscriberRepo interface {
	Insert(sub models.Subscriber) error
	Select(shopID string) ([]models.Subscriber, error)
	IsSubscribed(tgID, shopID string) bool
}

type subscriber struct {
	db *sqlx.DB
}

func (s subscriber) Insert(sub models.Subscriber) error {
	_, err := s.db.Query("insert into subscriber (username, first_name, last_name, telegram_id, avatar_url, shop_id) values ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING", sub.UserName, sub.FirstName, sub.LastName, sub.TelegramID, sub.AvatarUrl, sub.ShopID)
	return err
}

func (s subscriber) Select(shopID string) ([]models.Subscriber, error) {
	var subs []models.Subscriber
	err := s.db.Select(&subs, `select * from subscriber where shop_id = $1`, shopID)
	return subs, err
}

func (s subscriber) IsSubscribed(tgID, shopID string) bool {
	var id string
	s.db.Get(&id, "select id from subscriber where telegram_id = $1 and shop_id = $2", tgID, shopID)
	return id != ""
}
