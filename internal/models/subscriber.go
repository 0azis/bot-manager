package models

type Subscriber struct {
	ID          string `json:"ID" db:"id"`
	CreatedDate string `json:"createdDate" db:"created_date"`
	UpdateDate  string `json:"updatedDate" db:"updated_date"`
	FirstName   string `json:"firstName" db:"first_name"`
	LastName    string `json:"lastName" db:"last_name"`
	TelegramID  string `json:"telegramId" db:"telegram_id"`
	AvatarUrl   string `json:"avatarUrl" db:"avatar_url"`
	UserName    string `json:"userName" db:"username"`
	ShopID      string `json:"shopID" db:"shopId"`
}
