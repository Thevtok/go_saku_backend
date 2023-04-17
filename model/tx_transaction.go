package model

import "time"

type Transaction struct {
	Transaction_ID int       `json:"Transaction_ID"`
	User_ID        int       `json:"User_ID"`
	Type           string    `json:"Type"`
	Amount         float64   `json:"Amount"`
	Date           time.Time `json:"Date"`
	Description    string    `json:"Description"`
}
