package model

type PhotoUrl struct {
	Photo_ID uint   `json:"photo_id"`
	User_ID  uint   `json:"user_id"`
	Url      string `json:"url"`
}