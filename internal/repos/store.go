package repos 

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/stdlib"
)

type StoreRepos interface {
	Shop() shopRepo 
	Subscriber() subscriberRepo 
	Mail() mailRepo
}

type Store struct {
	db *sqlx.DB
}

func (s Store) Shop() shopRepo {
	return &shop{
		db: s.db,
	}	
}

func (s Store) Subscriber() subscriberRepo {
	return &subscriber{
		db: s.db,
	}
}

func (s Store) Mail() mailRepo {
	return &mail{
		db: s.db,
	}
}

func NewDB(user, password, host, port, dbName string) (Store, error) {
	// full URI using env variables
	DB_URI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	// getting db object
	db, err := sqlx.Connect("pgx", DB_URI)
	return Store{
		db: db,
	}, err
}
