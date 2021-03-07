package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	_ = setupTestDb()
	_ = testDb.AutoMigrate(&User{}, &Application{}, &Key{})
}

func TestKeyVerifier(t *testing.T) {

	var (
		key       Key
		parsedKey *Key
		claims    KeyClaims
	)

	key = Key{}
	key.ID = 42

	_, err := key.SetRandomSecret(32)

	testDb.Create(&key)
	testDb.Commit()

	assert.Nil(t, err, "key failed to WriteSecret")

	// normally client side builds our JWT here
	// JWT "claims" what key number it is

	claims = KeyClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 3).Unix(),
		},
		ID: key.ID,
	}

	clientToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	clientToken.Header["kid"] = claims.ID
	secretBytes, err := key.GetSecretBytes()

	assert.Nil(t, err, "GetSecretBytes() failed")

	ss, err := clientToken.SignedString(secretBytes)

	assert.Nil(t, err, "failed to sign JWT")

	verifier := KeyVerifier{ss, testDb}
	parsedKey, err = verifier.VerifyKey()

	assert.Nil(t, err, "verify failed")
	assert.Equal(t, parsedKey.ID, key.ID)
}
