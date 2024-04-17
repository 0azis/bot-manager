package models

import (
	"fmt"
	// "github.com/go-telegram/bot"
)

type Shop struct {
	ID     string `json:"ID" db:"id"`
	UserID string `json:"userID" db:"userId"`
	BotID  int    `json:"botID" db:"bot_id"`

	CreatedDate string `json:"createdDate" db:"createdDate"`
	UpdatedDate string `json:"updatedDate" db:"updatedDate"`

	Token       string `json:"token" db:"token"`
	TitleButton string `json:"titleButton" db:"titleButton"`
	Description string `json:"description" db:"description"`

	FirstName string `json:"firstName" db:"firstName"`
	Username  string `json:"username" db:"username"`

	Greetings   string `json:"greetings" db:"greetings"`
	FirstLaunch string `json:"firstLaunch" db:"firstLaunch"`
	AfterOrder  string `json:"afterOrder" db:"afterOrder"`

	IsActive bool `json:"isActive" db:"isActive"`
}

func WebLink(id string) string {
	var base = fmt.Sprintf("https://tgrocket.ru/shop/%s", id)
	return base
}

type ShopCredentials struct {
	Token string `json:"token"`
}

 