package repository

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ReygaFitra/inc-final-project.git/model"
)

func TestUserRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	columns := []string{"name", "username", "email", "phone_number", "address", "balance", "point", "tx_count"}
	rows := sqlmock.NewRows(columns).
		AddRow("John Doe", "johndoe", "johndoe@example.com", "+1-555-555-5555", "123 Main St", 1000, 10, 5).
		AddRow("Jane Smith", "janesmith", "janesmith@example.com", "+1-555-555-5555", "456 Main St", 500, 5, 3)
	mock.ExpectQuery("SELECT name,username,email,phone_number,address,balance ,point,tx_count from mst_users").WillReturnRows(rows)

	repo := &userRepository{db}

	users := repo.GetAll()
	if err != nil {
		t.Errorf("error was not expected while getting users: %s", err)
	}

	expectedUsers := []model.UserResponse{
		{Name: "John Doe", Username: "johndoe", Email: "johndoe@example.com", Phone_Number: "+1-555-555-5555", Address: "123 Main St", Balance: 1000, Point: 10, TxCount: 5},
		{Name: "Jane Smith", Username: "janesmith", Email: "janesmith@example.com", Phone_Number: "+1-555-555-5555", Address: "456 Main St", Balance: 500, Point: 5, TxCount: 3},
	}

	if !reflect.DeepEqual(expectedUsers, users) {
		t.Errorf("expected: %v but got: %v", expectedUsers, users)
	}
}
