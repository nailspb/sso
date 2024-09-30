package passwdHelper

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"sso/pkg/helpers/errorHelper"
)

const (
	ErrorCreateHash    = "Error on create password hash"
	ErrorWrongPassword = "Wrong password"
)

func HashPassword(password string) (string, error) {
	const op = "pkg.helpers.passwdHelper.HashPassword()"
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errorHelper.WrapError(op, ErrorCreateHash, err)
	}
	fmt.Println("password hashed", string(bytes))
	return string(bytes), nil
}

func ComparePassword(password, hash string) error {
	const op = "pkg.helpers.passwdHelper.ComparePassword()"
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errorHelper.WrapError(op, ErrorWrongPassword, err)
	}
	return nil
}
