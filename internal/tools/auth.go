package tools

import (
	"math/rand"
	acc "soa-hw-ilyaleshchyk/internal/account"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) ([]byte, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func IsSamePass(pass []byte, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, pass)
}

func GenerateRandomAccountID() acc.AccountID {
	return acc.AccountID(rand.Int31())
}
