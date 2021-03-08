package models

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type Key struct {
	gorm.Model
	ApplicationID uint // the application this key is valid for
	Enabled       bool
	Memo          *string
	HardwareID    *string
	Remaining     uint
	Secret        string
}

func (k *Key) SetRandomSecret(size int) (n int, err error) {
	b := make([]byte, size)
	n, err = rand.Read(b)

	if err != nil {
		return
	}

	k.Secret = base64.StdEncoding.EncodeToString(b)
	return
}

func (k *Key) GetSecretBytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(k.Secret)
}

func lookupKey(db *gorm.DB, keyID uint) ([]byte, error) {
	var key Key
	err := db.First(&key, keyID).Error

	if err != nil {
		return nil, err
	}

	return key.GetSecretBytes()
}

type KeyClaims struct {
	jwt.StandardClaims
	ApplicationID uint `json:"application_id"` // Key.ApplicationID field
	KeyID         uint `json:"key_id"`         // Key.ID Field
}

type KeyActivation struct {
	gorm.Model
	Key        Key
	Identifier *string
}
