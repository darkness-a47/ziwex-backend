package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHashedPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(pass), err
}

func CompareHashPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
