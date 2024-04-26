package database 

import (
	"botmanager/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type message struct {
	db *sqlx.DB
}

func (m message) Insert(msg domain.Message) error {
	_, err := m.db.Query(`insert into message (text, is_from_user, subscriber_id, bot_id) values ($1, $2, $3, $4)`, msg.Text, msg.IsFromUser, msg.SubscriberID, msg.BotID)
	return err
}

func (m message) Get(ID string) (domain.Message, error) {
	var message domain.Message
	err := m.db.Get(&message, `select * from message where id = $1`, ID)
	return message, err
}
