package auth

import (
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	passwordBytes := hasher.Sum(nil)

	bytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	passwordBytes := hasher.Sum(nil)

	err := bcrypt.CompareHashAndPassword([]byte(hash), passwordBytes)
	return err == nil
}
