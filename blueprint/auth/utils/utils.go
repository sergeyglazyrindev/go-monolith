package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// bcryptDiff
var bcryptDiff = 12

// HashPass Generates a hash from a password and salt
func HashPass(pass string, salt string) (string, error) {
	password := []byte(pass + salt)
	hash, err := bcrypt.GenerateFromPassword(password, bcryptDiff)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
