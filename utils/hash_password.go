package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Genereate hash passowrd from string
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", errors.New("password hash error")
	}

	return string(hash), nil
}

// Verify password is valid
func VerifyPassword(hash string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
