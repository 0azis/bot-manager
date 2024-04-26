package domain

type User struct {
	ID           string `db:"id"`
	Code         string `db:"code"`
	TelegramID   string `db:"telegram_id"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	UserName     string `db:"user_name"`
	IsPremium    bool   `db:"is_premium"`
	LanguageCode string `db:"language_code"`
	AvatarUrl    string `db:"avatar_url"`
	CreatedDate  string `db:"created_date"`
	UpdatedDate  string `db:"updated_date"`
}
