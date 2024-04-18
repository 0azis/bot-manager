package models

import (
	"fmt"
	// "github.com/go-telegram/bot"
)

type Shop struct {
	ID     string `db:"id"`
	UserID string `db:"userId"`
	BotID  int    `db:"bot_id"`

	CreatedDate string `db:"createdDate"`
	UpdatedDate string `db:"updatedDate"`

	Token       string `db:"token"`
	TitleButton string `db:"titleButton"`
	Description string `db:"description"`

	FirstName string `db:"firstName"`
	Username  string `db:"username"`

	Greetings   string `db:"greetings"`
	FirstLaunch string `db:"firstLaunch"`
	AfterOrder  string `db:"afterOrder"`

	IsActive bool `db:"isActive"`
}

func WebLink(id string) string {
	var base = fmt.Sprintf("https://tgrocket.ru/shop/%s", id)
	return base
}

type ShopCredentials struct {
	Token string `json:"token"`
}
