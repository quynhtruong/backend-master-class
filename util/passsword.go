package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hassedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password %w", err)
	}
	return string(hassedPassword), nil
}

func CheckPassword(password string, hassedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hassedPassword), []byte(password))
}
