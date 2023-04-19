package model

type Reward struct {
	Point_ID      uint `json:"point_id"`
	User_ID       int  `json:"user_id"`
	Amount_Reward int  `json:"amount_reward"`
}

type UserPoint struct {
	Point_ID      uint   `json:"point_id"`
	User_ID       int    `json:"user_id"`
	Amount_Reward int    `json:"amount_reward"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone_Number  string `json:"phone_number"`
	Balance       int    `json:"balance"`
}