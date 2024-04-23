package repo

import (
	"botmanager/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type subscriber struct {
	db *sqlx.DB
}

func (s subscriber) Insert(sub domain.Subscriber) error {
	_, err := s.db.Query("insert into subscriber (username, first_name, last_name, telegram_id, avatar_url, shop_id) values ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING", sub.UserName, sub.FirstName, sub.LastName, sub.TelegramID, sub.AvatarUrl, sub.ShopID)
	return err
}

func (s subscriber) Select(shopID string) ([]domain.Subscriber, error) {
	var subs []domain.Subscriber
	err := s.db.Select(&subs, `select * from subscriber where shop_id = $1`, shopID)
	return subs, err
}

func (s subscriber) IsSubscribed(tgID, shopID string) bool {
	var id string
	s.db.Get(&id, "select id from subscriber where telegram_id = $1 and shop_id = $2", tgID, shopID)
	return id != ""
}

func (s subscriber) Get(ID string) (domain.Subscriber, error) {
	var subscriber domain.Subscriber
	err := s.db.Get(&subscriber, `select * from subscriber where id = $1`, ID)
	return subscriber, err
}

func (s subscriber) GetByTelegramID(telegramID string) (domain.Subscriber, error) {
	var subscriber domain.Subscriber
	err := s.db.Get(&subscriber, `select * from subscriber where telegram_id = $1`, telegramID)
	return subscriber, err
}
