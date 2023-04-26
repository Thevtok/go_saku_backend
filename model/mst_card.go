package model

type Card struct {
	CardID         uint   `json:"card_id"`
	UserID         uint   `json:"user_id"`
	CardType       string `json:"card_type"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
	Name           string `json:"name"`
}

type CardResponse struct {
	UserID         uint   `json:"user_id"`
	CardType       string `json:"card_type"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}
