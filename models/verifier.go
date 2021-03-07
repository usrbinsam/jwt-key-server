package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type KeyVerifierConfig struct {
	Secret []byte
}

type KeyVerifier struct {
	JWT string
	DB  *gorm.DB
}

func (v *KeyVerifier) VerifyKey() (*Key, error) {
	key := Key{}
	err := key.FromJWT(v.JWT, func(token *jwt.Token) (interface{}, error) {
		f, ok := token.Header["kid"]

		if !ok {
			return nil, errors.New("kid not set on Header")
		}

		kid, ok := f.(float64)

		if !ok {
			return nil, errors.New("kid on Header is invalid data type")
		}

		return lookupKey(v.DB, uint(kid))
	})

	if err != nil {
		return nil, err
	}

	return &key, nil
}
