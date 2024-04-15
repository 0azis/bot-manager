package repos

import (
	"botmanager/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// interface for shop
type ShopRepo interface {
	Select() ([]models.Shop, error)
	Get(token string) (models.Shop, error)
	UpdateStatus(token string, status bool) error
}

type shop struct {
	db *sqlx.DB
}

// get all shops from DB
func (s shop) Select() ([]models.Shop, error) {
	shops := []models.Shop{}
	err := s.db.Select(&shops, `select * from shop`)
	return shops, err
}

// get one shop from DB
func (s shop) Get(token string) (models.Shop, error) {
	shop := models.Shop{}
	err := s.db.Get(&shop, `select * from shop where token = $1`, token)
	return shop, err
}

// update work status of shop
func (s shop) UpdateStatus(token string, status bool) error {
	fmt.Println(fmt.Sprintf("update shop set isActive = %t where token = '%s'", status, token))
	_, err := s.db.Query(fmt.Sprintf("update shop set isActive = %t where token = '%s'", status, token))
	return err
}

// init a new shop repository
func NewShopRepo(db *sqlx.DB) *shop {
	return &shop{
		db: db,
	}
}
