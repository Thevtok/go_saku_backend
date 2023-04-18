package model

type Reward struct {
	Point_ID      uint `json:"point_id"`
	User_ID       int  `json:"user_id"`
	Amount_Reward int  `json:"amount_reward"`
}