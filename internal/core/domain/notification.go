package domain

type Notification struct {
	UserID   string `json:"userID"`
	ShopID   string `json:"shopID"`
	Text     string `json:"text"`
	Photo    string `json:"photo"`
	Button   bool   `json:"button"`
	Link     string `json:"link"`
	LinkText string `json:"linkText"`
}
