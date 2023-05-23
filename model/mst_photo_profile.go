package model

type PhotoUrl struct {
	Photo_ID uint   `json:"photo_id"`
	UserID   string `json:"user_id"`
	Url      string `json:"url"`
}
