package services

import (
	"crypto/rsa"
	"sso/pkg/helpers/errorHelper"
	"sso/pkg/helpers/jwtHelper"
)

type Storage interface {
	GetRsaKey() (*rsa.PrivateKey, error)
}

func Check(token string, storage Storage) error {
	const op = "internal.services.check"
	key, _ := storage.GetRsaKey()
	err := jwtHelper.Check(&key.PublicKey, token)
	if err != nil {
		return errorHelper.WrapError(op, "token is invalid", err)
	}
	return nil
}
