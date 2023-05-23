package repository

import (
	"database/sql"
	"fmt"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

type PhotoRepository interface {
	Create(photo *model.PhotoUrl) error
	GetByID(id string) (*model.PhotoUrl, error)
	Update(photo *model.PhotoUrl) error
	Delete(id string) string
}

type photoRepository struct {
	db *sql.DB
}

func (r *photoRepository) Create(photo *model.PhotoUrl) error {
	_, err := r.db.Exec("INSERT INTO mst_photo_url (url_photo, user_id) VALUES ($1, $2)", photo.Url, photo.UserID)
	if err != nil {
		return fmt.Errorf("failed to create photo: %w", err)
	}
	return nil
}

func (r *photoRepository) GetByID(id string) (*model.PhotoUrl, error) {
	var user model.PhotoUrl
	row := r.db.QueryRow("SELECT url_photo, user_id from mst_photo_url WHERE user_id = $1", id)
	err := row.Scan(&user.Url, &user.UserID)
	if err != nil {
		return nil, err // Mengembalikan kesalahan jika terjadi error saat pemindaian
	}
	return &user, nil
}

func (r *photoRepository) Update(photo *model.PhotoUrl) error {
	query := "UPDATE mst_photo_url SET url_photo = $1 WHERE user_id = $2"
	_, err := r.db.Exec(query, photo.Url, photo.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *photoRepository) Delete(id string) string {
	query := "DELETE FROM mst_photo_url WHERE user_id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Println("failed to delete photo")
	}
	return "Delete photo successfully"
}
func NewPhotoRepository(db *sql.DB) PhotoRepository {
	return &photoRepository{db: db}
}
