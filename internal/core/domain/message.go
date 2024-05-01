package domain

type Message struct {
	ID           string `json:"ID" db:"id"`
	CreatedDate  string `json:"from" db:"created_date"`
	Text         string `json:"text" db:"text"`
	IsFromUser   bool   `json:"isFromUser" db:"is_from_user"`
	SubscriberID string `json:"subscrbierID" db:"subscriber_id"`
	BotID        string `json:"botID" db:"bot_id"`
}

type MessageCredentials struct {
	ID string `json:"ID"`
}
