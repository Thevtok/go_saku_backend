package model

type Card struct {
	Card_ID         int    `json:"Card_ID"`
	User_ID         int    `json:"User_ID"`
	Card_Type       string `json:"Card_Type"`
	Card_Number     string `json:"Card_Number"`
	Expiration_Date string `json:"Expiration_Date"`
	CVV             string `json:"CVV"`
}
