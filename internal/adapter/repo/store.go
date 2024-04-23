package repo

import (
	"botmanager/internal/core/port/repository"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/stdlib"
)

// Store struct with all repos
type Store struct {
	Shop       repository.ShopRepository
	Subscriber repository.SubscriberRepository
	Message    repository.MessageRepository
	Mail       repository.MailRepository
}

func NewDB(user, password, host, port, dbName string) (*Store, error) {
	// full URI using env variables
	DB_URI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	// getting db object
	db, err := sqlx.Connect("pgx", DB_URI)
	s := &Store{
		Shop:       shop{db},
		Subscriber: subscriber{db},
		Message:    message{db},
		Mail:       mail{db},
	}
	return s, err
}
