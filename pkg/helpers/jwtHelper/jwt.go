package jwtHelper

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	errorHelper "sso/pkg/helpers/errorHelper"
)

const (
	ErrorCreateToken  = "errorHelper creating the token"
	ErrorCheckSign    = "errorHelper on checked the signature"
	ErrorInvalidToken = "token is invalid"
	ErrorGetClaim     = "errorHelper on get claim"
	ErrorParseToken   = "errorHelper on parse token"
)

func Create(key *rsa.PrivateKey, claim map[string]any) (string, error) {
	const op = "pkg.helpers.jwtHelper.create()"
	j := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claim))
	if t, err := j.SignedString(key); err != nil {
		return "", errorHelper.WrapError(op, ErrorCreateToken, err)
	} else {
		return t, nil
	}
}

func GetClaim(key *rsa.PublicKey, tokenString string) (*jwt.MapClaims, error) {
	const op = "pkg.helpers.jwtHelper.getClaim()"
	if token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		const op = "pkg.helpers.jwtHelper.getKeyFunc()"
		// Проверка, что метод подписи RSA:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return false, errorHelper.WrapError(op, ErrorCheckSign, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		// возвращаем указатель на публичный ключ
		return key, nil

	}); err != nil {
		return nil, errorHelper.WrapError(op, ErrorParseToken, err)
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			return &claims, nil
		} else {
			return nil, errorHelper.WrapError(op, ErrorGetClaim, err)
		}
	}
}

func Check(key *rsa.PublicKey, tokenString string) error {
	const op = "pkg.helpers.jwtHelper.check()"
	if _, err := GetClaim(key, tokenString); err != nil {
		return errorHelper.WrapError(op, ErrorInvalidToken, err)
	}
	return nil
}
