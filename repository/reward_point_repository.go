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
	// Update(point *model.Reward) string
	Delete(id uint) string
}

type rewardRepository struct {
	db *sql.DB
}

func (r *rewardRepository) GetAll() any {
	var users []model.UserPoint
	query := "SELECT u.name, u.email, u.phone_number, u.balance, r.amount_reward FROM mst_reward_point AS r INNER JOIN mst_users AS u ON r.user_id = u.user_id"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user model.UserPoint
		if err := rows.Scan(&user.Name, &user.Email, &user.Phone_Number, &user.Balance, &user.Amount_Reward); err != nil {
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
	var pointDB model.Reward
	var userDB model.User
	query := "SELECT u.name, u.email, u.phone_number, u.balance, r.amount_reward FROM mst_reward_point AS r INNER JOIN mst_users AS u ON r.user_id = u.user_id WHERE r.user_id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&pointDB.User_ID, &userDB.Name, &userDB.Email, &userDB.Phone_Number, &userDB.Balance, &pointDB.Amount_Reward)
	if err != nil {
		log.Println(err)
	}
	if userDB.ID == 0 {
		return "user not found!"
	}
	return pointDB
}

func (r *rewardRepository) Create(newPoint *model.Reward) string {
	stmt, err := r.db.Prepare("INSERT INTO mst_reward_point (user_id, amount_reward) VALUES ($1, $2)")
	if err != nil {
		return err.Error()
	}
	defer stmt.Close()
	_, err = stmt.Exec(newPoint.User_ID, newPoint.Amount_Reward)
	if err != nil {
		return err.Error()
	}
	return "Point Added!"
}

// func (r *rewardRepository) Update(point *model.Reward) string {
// 	res := r.GetByID(point.Point_ID)
// 	if res == "user not found!" {
// 		return res.(string)
// 	}

// 	query := "UPDATE mst_reward_point SET user_id = $1, amount_reward = $2 WHERE point_id = $1"
// 	_, err := r.db.Exec(query, point.User_ID, point.Amount_Reward, point.Point_ID)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return fmt.Sprintf("user point with id %d updated successfully", point.Point_ID)
// }

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
