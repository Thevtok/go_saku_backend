package model

type User struct {
	ID           string `json:"user_id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone_Number string `json:"phone_number"`
	Address      string `json:"address"`
	Balance      int    `json:"balance"`
	Role         string `json:"role"`
	Point        int    `json:"point"`
	Token        string `json:"token"`
	Badge        string `json:"badge_name"`
	BadgeID      int    `json:"badge_id"`
	TxCount      int    `json:"tx_count"`
}

// Define the table name for the User struct
