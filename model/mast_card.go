package model

type Card struct {
	CardID         uint   `json:"card_id"`
	Username       string `json:"username"`
	CardType       string `json:"card_type"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
	Name           string `json:"name"`
}

type CardResponse struct {
	Username       string `json:"username"`
	CardType       string `json:"card_type"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}
