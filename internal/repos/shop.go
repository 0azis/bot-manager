package repos

import (
	"botmanager/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// interface for shop
type shopRepo interface {
	Select() ([]models.Shop, error)
	Get(token string) (models.Shop, error)
	UpdateStatus(token string, status bool) error
	IsTokenValid(token string) bool
	GetBy(key string, value any) (models.Shop, error)
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

func (s shop) GetBy(key string, value any) (models.Shop, error) {
	var shopID models.Shop 
	// fmt.Println(fmt.Sprintf("select * from shop where %s = '%s'", key, value))
	err := s.db.Get(&shopID, fmt.Sprintf("select * from shop where %s = '%s'", key, value))
	return shopID, err
}

// update work status of shop
func (s shop) UpdateStatus(token string, status bool) error {
	_, err := s.db.Query(fmt.Sprintf("update shop set isActive = %t where token = '%s'", status, token))
	return err
}

func (s shop) IsTokenValid(token string) bool {
	var isToken string
	err := s.db.Get(&isToken, `select token from shop where token = $1`, token)
	return err != nil
}
