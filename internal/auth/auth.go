package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(pass string) (string, error) {
	passHash, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return passHash, nil
}

func CheckPasswordHash(pass, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(pass, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
