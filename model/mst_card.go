package model

type Card struct {
	CardID         uint   `json:"card_id"`
	UserID         uint   `json:"user_id"`
	CardType       string `json:"card_type"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}

type CardResponse struct {
	UserID         uint   `json:"user_id"`
	CardID         uint   `json:"card_id"`
	CardType       string `json:"card_type"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}

type CreateCard struct {
	UserID         uint   `json:"user_id"`
	CardType       string `json:"card_type"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}
