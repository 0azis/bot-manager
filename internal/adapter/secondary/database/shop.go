package database

import (
	"botmanager/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type shop struct {
	db *sqlx.DB
}

// get all shops from DB
func (s shop) Select() ([]domain.Shop, error) {
	shops := []domain.Shop{}
	err := s.db.Select(&shops, `select * from shop`)
	return shops, err
}

// get one shop from DB
func (s shop) Get(ID string) (domain.Shop, error) {
	shop := domain.Shop{}
	err := s.db.Get(&shop, `select * from shop where id = $1`, ID)
	return shop, err
}

func (s shop) GetByBotID(botID string) (domain.Shop, error) {
	shop := domain.Shop{}
	err := s.db.Get(&shop, `select * from shop where bot_id = $1`, botID)
	return shop, err
}

func (s shop) GetByToken(token string) (domain.Shop, error) {
	shop := domain.Shop{}
	err := s.db.Get(&shop, `select * from shop where token = $1`, token)
	return shop, err
}
