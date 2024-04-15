package setup

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/stdlib"
)

func NewDB(user, password, host, port, dbName string) (*sqlx.DB, error) {
	// full URI using env variables
	DB_URI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	// getting db object
	db, err := sqlx.Connect("pgx", DB_URI)
	return db, err
}
