package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

type PhotoRepository interface {
	Create(photo *model.PhotoUrl) error
	GetByID(id uint) (*model.PhotoUrl, error)
	Update(photo *model.PhotoUrl) error
	Delete(id uint) string
}

type photoRepository struct {
	db *sql.DB
}

func NewPhotoRepository(db *sql.DB) PhotoRepository {
	return &photoRepository{db: db}
}

func (r *photoRepository) Create(photo *model.PhotoUrl) error {
    _, err := r.db.Exec("INSERT INTO mst_photo_url (user_id, url_photo) VALUES ($1, $2)", photo.User_ID, photo.Url)
    if err != nil {
        return err
    }
    return nil
}

func (r *photoRepository) GetByID(id uint) (*model.PhotoUrl, error) {
	var user model.PhotoUrl
	row := r.db.QueryRow("SELECT user_id, url_photo from mst_photo_url WHERE user_id = $1", id)
	err := row.Scan(&user.User_ID, &user.Url)
	if err != nil {
		log.Println(err)
	}
	return &user, nil
}

func (r *photoRepository) Update(photo *model.PhotoUrl) error {
	query := "UPDATE mst_photo_url SET url_photo=$1 WHERE user_id=$2"
	_, err := r.db.Exec(query, photo.Url, photo.User_ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *photoRepository) Delete(id uint) string {
	query := "DELETE FROM mst_photo_url WHERE user_id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Println("failed to delete photo")
	}
	return "Delete photo successfully"
}