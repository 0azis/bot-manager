package repo

import (
	"botmanager/internal/core/domain"
	// "fmt"

	"github.com/jmoiron/sqlx"
)

type user struct {
	db *sqlx.DB
}

func (u user) Insert(user domain.User) error {
	// fmt.Println(user)
	// fmt.Println(fmt.Sprintf(`insert into "user" (code, telegram_id, first_name, last_name, user_name, is_premium, language_code, avatar_url) values (%d, '%s', '%s', '%s', '%s', %t, '%s', '%s') on conflict (telegram_id) do update set first_name='%s', last_name='%s', user_name='%s', is_premium=%t, language_code='%s', avatar_url='%s';`, user.Code, user.TelegramID, user.FirstName, user.LastName, user.UserName, user.IsPremium, user.LanguageCode, user.AvatarUrl, user.FirstName, user.LastName, user.UserName, user.IsPremium, user.LanguageCode, user.AvatarUrl))
	_, err := u.db.Query(`insert into "user" (code, telegram_id, first_name, last_name, user_name, is_premium, language_code, avatar_url) values ($1, $2, $3, $4, $5, $6, $7, $8) on conflict (telegram_id) do update set first_name=$3, last_name=$4, user_name=$5, is_premium=$6, language_code=$7, avatar_url=$8;`, user.Code, user.TelegramID, user.FirstName, user.LastName, user.UserName, user.IsPremium, user.LanguageCode, user.AvatarUrl)
	return err
}
