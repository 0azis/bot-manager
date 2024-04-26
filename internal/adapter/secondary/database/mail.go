package database 

import (
	"botmanager/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type mail struct {
	db *sqlx.DB
}

func (m mail) Get(mailID string) (domain.Mail, error) {
	var mail domain.Mail
	err := m.db.Get(&mail, `select * from share where id = $1`, mailID)
	return mail, err
}
