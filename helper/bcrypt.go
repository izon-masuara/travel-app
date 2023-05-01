package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(plaintext string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), 14)
	return string(bytes), err
}

func DecribePassword(passowrd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passowrd))
	return err == nil
}
