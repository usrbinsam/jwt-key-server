package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// KeyVerifier is responsible for verifying the authenticity of a given JWT and the validity of the claims within the
// JWT. This interface should be used for verifying end-user key operations such as activation and verification.
type KeyVerifier struct {
	JWT string
	DB  *gorm.DB
}

// VerifyKey ensures that both the JWT is signed and the claims within the JWT are valid. Use this function to ensure
// that a key is trusted and belongs to the Application trying to activate with it.
func (v *KeyVerifier) VerifyKey() (*Key, error) {
	claims, err := v.ExtractClaims(v.keyFunc)

	if err != nil {
		return nil, err
	}

	key := Key{}
	err = v.DB.Take(&key, claims.KeyID).Error

	if err != nil {
		return nil, err
	}

	if key.ApplicationID != claims.ApplicationID {
		return nil, errors.New("this key is for another application")
	}

	return &key, nil
}

func (v *KeyVerifier) keyFunc(token *jwt.Token) (interface{}, error) {
	f, ok := token.Header["kid"]

	if !ok {
		return nil, errors.New("kid not set on Header")
	}

	kid, ok := f.(float64)

	if !ok {
		return nil, errors.New("kid on Header is invalid data type")
	}

	return lookupKey(v.DB, uint(kid))
}

// ExtractClaims() loads verified claims from the JWT field.
func (v *KeyVerifier) ExtractClaims(keyLookup jwt.Keyfunc) (*KeyClaims, error) {

	token, err := jwt.ParseWithClaims(v.JWT, &KeyClaims{}, keyLookup)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		panic("invalid token, but no error returned from jwt.ParseWithClaims")
	}

	claims, ok := token.Claims.(*KeyClaims)
	if !ok {
		return nil, errors.New("invalid claims structure")
	}
	return claims, nil
}
