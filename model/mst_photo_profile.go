package model

type PhotoUrl struct {
	Photo_ID uint   `json:"photo_id"`
	Url      string `json:"url"`
	Username string `json:"username"`
}