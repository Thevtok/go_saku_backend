package model

type PhotoUrl struct {
	Photo_ID uint   `json:"photo_id"`
	UserID   uint   `json:"user_id"`
	Url      string `json:"url"`
}
