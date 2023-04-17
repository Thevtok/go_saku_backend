package model

type Users struct {
	User_ID      int     `json:"User_ID"`
	Name         string  `json:"Name"`
	Email        string  `json:"Email"`
	Password     string  `json:"Password"`
	Phone_Number string  `json:"Phone_Number"`
	Address      string  `json:"Address"`
	Balance      float64 `json:"Balance"`
}
