package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HasingPassword(password string) (string, error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return "", err
	}
	return string(hashedByte), err
}

func CheckPasswordHash(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("password hash does not match")
	}
	return nil
}
