package models

import (
	"errors"
	"fmt"
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

		if kid, ok := token.Header["kid"]; ok {
			var id uint

			switch v := kid.(type) {
			case float64:
				id = uint(token.Header["kid"].(float64))
			case int:
				id = uint(token.Header["kid"].(int))
			default:
				return nil, errors.New(fmt.Sprintf("kid header incorrect type %T", v))
			}

			return lookupKey(v.DB, id)
		}
		return nil, errors.New("kid not in header")
	})
	if err != nil {
		return nil, err
	}
	return &key, nil
}
