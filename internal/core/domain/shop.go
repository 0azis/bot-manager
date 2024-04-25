package domain

import (
	"fmt"
)

type Shop struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
	BotID  int    `db:"bot_id"`

	CreatedDate string `db:"created_date"`
	UpdatedDate string `db:"updated_date"`

	Token       string `db:"token"`
	TitleButton string `db:"title_button"`
	Description string `db:"description"`

	FirstName string `db:"first_name"`
	Username  string `db:"username"`

	Greetings   string `db:"greetings"`
	FirstLaunch string `db:"first_launch"`
	AfterOrder  string `db:"after_order"`

	IsActive bool `db:"is_active"`
}

func WebLink(id string) string {
	var base = fmt.Sprintf("https://tgrocket.ru/shop/%s", id)
	return base
}

type ShopCredentials struct {
	Token string `json:"token"`
}

type ShopBot struct {
	Data Shop
}

type HomeBot struct {
	Token string
}
