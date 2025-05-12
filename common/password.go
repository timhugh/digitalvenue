package common

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password, salt string) (string, string, error) {
	if salt == "" {
		salt = uuid.NewString() // TODO: uuid isn't right, generate actual salt
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(salt+password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return string(hashedPassword), salt, nil
}
