package domain 

type Mail struct {
	ID          string `db:"id"`
	CreatedDate string `db:"createdDate"`
	UpdatedDate string `db:"updatedDate"`
	Text        string `db:"text"`
	AddButton   bool   `db:"addButton"`
	PhotoLink   string `db:"photoLink"`
	ShopID      string `db:"shop_id"`
}

type MailCredentials struct {
	ID string `json:"ID"`
}
