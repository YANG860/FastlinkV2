package api

import (
	"golang.org/x/crypto/bcrypt"
)

func hash(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashed)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isValidUsername(username string) bool {
	// Add your username validation logic here
	return len(username) >= 3 && len(username) <= 20
}

func isValidPassword(password string) bool {
	// Add your password validation logic here
	return len(password) >= 6
}
