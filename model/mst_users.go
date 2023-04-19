package model

type User struct {
	ID           uint   `gorm:"column:user_id;primary_key"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone_Number string `json:"phone_number"`
	Address      string `json:"address"`
	Balance      int    `json:"balance"`
}

// Define the table name for the User struct
type UserGetAll struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone_Number string `json:"phone_number"`
	Address      string `json:"address"`
	Balance      int    `json:"balance"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserID   int    `json:"user_id"`
}
