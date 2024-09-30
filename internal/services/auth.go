package services

import (
	"errors"
	"sso/internal/storage"
	"sso/pkg/helpers/errorHelper"
	"sso/pkg/helpers/jwtHelper"
	"sso/pkg/helpers/passwdHelper"
	"time"
)

const (
	ErrorQueryUser    = "failed request for user data"
	ErrorGetRsaKey    = "failed get rsa key"
	ErrorUserNotFound = "user not found (bad login)"
	ErrorCreateToken  = "failed to create token"
)

func Auth(login string, password string, storage storage.Storage) (string, error) {
	const op = "internal.services.auth"
	if u, err := storage.Users().GetUser(login); err != nil {
		return "", errorHelper.WrapError(op, ErrorQueryUser, err)
	} else if u != nil {
		if err := passwdHelper.ComparePassword(password, u.Password); err == nil {
			key, err := storage.GetRsaKey()
			if err != nil {
				return "", errorHelper.WrapError(op, ErrorQueryUser, err)
			}
			if token, err := jwtHelper.Create(key, map[string]any{
				"iss": "DM SSO",
				"sub": "auth",
				"aud": u.Id,
				"exp": time.Now().Add(time.Hour * 2).Unix(),
				"nbf": time.Now().Unix(),
				"iat": time.Now().Unix(),
				"jti": "none", //token id
			}); err != nil {
				return "", errorHelper.WrapError(op, ErrorCreateToken, err)
			} else {
				return token, nil
			}
		} else {
			return "", errorHelper.WrapError(op, ErrorCreateToken, err)
		}
	} else {
		return "", errorHelper.WrapError(op, ErrorCreateToken, errors.New(ErrorUserNotFound))
	}
}
