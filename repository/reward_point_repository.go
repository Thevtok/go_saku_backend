package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"

	_ "github.com/lib/pq"
)

type RewardRepository interface {
	GetAll() any
	GetByID(id uint) any
	Create(newPoint *model.Reward) string
	Update(point *model.Reward) string
	Delete(id uint) string
}

type rewardRepository struct {
	db *sql.DB
}

func (r *rewardRepository) GetAll() any {
	var users []model.Reward
	query := "SELECT point_id, user_id, amount_reward FROM mst_reward_point"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user model.Reward
		if err := rows.Scan(&user.Point_ID, &user.User_ID, &user.Amount_Reward); err != nil {
			log.Println(err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	if len(users) == 0 {
		return "no data"
	}
	return users
}

func (r *rewardRepository) GetByID(id uint) any {
	var userDB model.Reward
	query := "SELECT amount_reward FROM mst_reward_point WHERE user_id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&userDB.Point_ID, &userDB.User_ID, &userDB.Amount_Reward)
	if err != nil {
		log.Println(err)
	}
	if userDB.Point_ID == 0 {
		return "user not found!"
	}
	return userDB
}

func (r *rewardRepository) Create(newPoint *model.Reward) string {
	query := "INSERT INTO mst_reward_point (user_id, amount_reward) VALUES ($1, $2)"
	_, err := r.db.Exec(query, newPoint.User_ID, newPoint.Amount_Reward)

	if err != nil {
		log.Println(err)
		return "failed to add point!"
	}
	return "point create succesfully"
}

func (r *rewardRepository) Update(point *model.Reward) string {
	res := r.GetByID(point.Point_ID)
	if res == "user not found!" {
		return res.(string)
	}

	query := "UPDATE mst_reward_point SET user_id = $1, amount_reward = $2 WHERE point_id = $1"
	_, err := r.db.Exec(query, point.User_ID, point.Amount_Reward, point.Point_ID)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("user point with id %d updated successfully", point.Point_ID)
}

func (r *rewardRepository) Delete(id uint) string {
	res := r.GetByID(id)
	if res == "user point not found!" {
		return res.(string)
	}
	query := "DELETE FROM mst_reward_point WHERE point_id = $1"
	_, err := r.db.Exec(query, id)

	if err != nil {
		log.Println(err)
		return "failed to delete user point"
	}
	return fmt.Sprintf("user point with id %d deleted successfully", id)
}

func NewRewardRepo(db *sql.DB) RewardRepository {
	repo := new(rewardRepository)
	repo.db = db
	return repo
}
