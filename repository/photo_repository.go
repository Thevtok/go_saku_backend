package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

type PhotoRepository interface {
	Create(photo *model.PhotoUrl) error
	GetByID(username string) (*model.PhotoUrl, error)
	Update(photo *model.PhotoUrl) error
	Delete(username string) string
}

type photoRepository struct {
	db *sql.DB
}

func NewPhotoRepository(db *sql.DB) PhotoRepository {
	return &photoRepository{db: db}
}

func (r *photoRepository) Create(photo *model.PhotoUrl) error {
    _, err := r.db.Exec("INSERT INTO mst_photo_url (url_photo, username) VALUES ($1, $2)", photo.Url, photo.Username)
    if err != nil {
        return err
    }
    return nil
}

func (r *photoRepository) GetByID(username string) (*model.PhotoUrl, error) {
	var user model.PhotoUrl
	row := r.db.QueryRow("SELECT url_photo, username from mst_photo_url WHERE username = $1", username)
	err := row.Scan(&user.Url, &user.Username)
	if err != nil {
		log.Println(err)
	}
	return &user, nil
}

func (r *photoRepository) Update(photo *model.PhotoUrl) error {
	query := "UPDATE mst_photo_url SET url_photo=$1 WHERE username=$2"
	_, err := r.db.Exec(query, photo.Url, photo.Username)
	if err != nil {
		return err
	}
	return nil
}

func (r *photoRepository) Delete(username string) string {
	query := "DELETE FROM mst_photo_url WHERE username = $1"
	_, err := r.db.Exec(query, username)
	if err != nil {
		fmt.Println("failed to delete photo")
	}
	return "Delete photo successfully"
}