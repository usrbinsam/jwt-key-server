package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type KeyStatus = uint8

const (
	KeyStatusActive = KeyStatus(1)
)

type Key struct {
	gorm.Model
	ApplicationID  uint // the application this key is valid for
	Status         KeyStatus
	Memo           *string
	Secret         string
	MaxActivations uint
	Activations    []*KeyActivation
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

func (k *Key) Activate(ka *KeyActivation) error {
	if k.Status != KeyStatusActive {
		return errors.New("key is not active")
	}

	if len(k.Activations) >= int(k.MaxActivations) {
		return errors.New("key has exceeded activations")
	}

	for _, activation := range k.Activations {
		if activation.Identifier == ka.Identifier {
			return errors.New("device is already activated")
		}
	}

	ka.KeyID = k.ID
	k.Activations = append(k.Activations, ka)

	return nil
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
	KeyID      uint
	Identifier *string `gorm:"uniqueIndex"`
}
