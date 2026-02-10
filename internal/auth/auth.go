package auth

import (
	"github.com/alexedwards/argon2id"
	"errors"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("Password cannot be empty")
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	if password == "" || hash == "" {
		return false, errors.New("Password and/or hash cannot be empty")
	}

	isMatch, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return isMatch, nil
}