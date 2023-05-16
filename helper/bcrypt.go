package helper

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plaintext string) (string, error) {
	cost := os.Getenv("BYCRYPT_COST")
	num, err := strconv.Atoi(cost)
	if err != nil {
		return "", err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), num)
	return string(bytes), err
}

func DecribePassword(passowrd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passowrd))
	return err == nil
}
