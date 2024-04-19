package repos

import (
	"botmanager/internal/models"

	"github.com/jmoiron/sqlx"
)

type mailRepo interface {
	Get(mailID string) (models.Mail, error)
}

type mail struct {
	db *sqlx.DB
}

func (m mail) Get(mailID string) (models.Mail, error) {
	var mail models.Mail
	err := m.db.Get(&mail, `select * from share where id = $1`, mailID)
	return mail, err
}