package database

import (
	"botmanager/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type user struct {
	db *sqlx.DB
}

// Add user after registration
func (u user) Insert(user domain.User) error {
	_, err := u.db.Query(`insert into "user" (telegram_id, first_name, last_name, user_name, is_premium, language_code, avatar_url) values ($1, $2, $3, $4, $5, $6, $7) on conflict (telegram_id) do update set first_name=$2, last_name=$3, user_name=$4, is_premium=$5, language_code=$6, avatar_url=$7;`, user.TelegramID, user.FirstName, user.LastName, user.UserName, user.IsPremium, user.LanguageCode, user.AvatarUrl)
	return err
}

func (u user) Get(ID string) (domain.User, error) {
	var user domain.User
	err := u.db.Get(&user, `select * from "user" where id = $1`, ID)
	return user, err
}
