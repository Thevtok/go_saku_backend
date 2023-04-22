package model

type User struct {
	ID           uint   `json:"user_id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone_Number string `json:"phone_number"`
	Address      string `json:"address"`
	Balance      int    `json:"balance"`
	Role         string `json:"role"`
	Point        int    `json:"point"`
}

// Define the table name for the User struct
type UserResponse struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phone_Number string `json:"phone_number"`
	Address      string `json:"address"`
	Balance      int    `json:"balance"`
	Point        int    `json:"point"`
}
type UserCreate struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone_Number string `json:"phone_number"`
	Address      string `json:"address"`
	Balance      int    `json:"balance"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
