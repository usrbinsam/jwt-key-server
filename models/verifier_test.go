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

	_, err := key.SetRandomSecret(64)

	testDb.Create(&key)
	testDb.Commit()

	assert.Nil(t, err, "key failed to SetRandomSecret")

	// normally client side builds our JWT here
	// JWT "claims" what key number it is

	expiration := time.Now().Add(time.Second * 3).Unix()
	claims = KeyClaims{
		ID: key.ID,
	}

	claims.ExpiresAt = expiration
	clientToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	clientToken.Header["kid"] = claims.ID

	// sign JWT with generated secret
	// perhaps application loads this key from a file or Windows Registry
	secretBytes, _ := key.GetSecretBytes()

	t.Run("TestValidJWT", func(t *testing.T) {
		ss, _ := clientToken.SignedString(secretBytes)

		verifier := KeyVerifier{ss, testDb}
		parsedKey, err = verifier.VerifyKey()

		assert.Nil(t, err, "verify failed")
		assert.Equal(t, parsedKey.ID, key.ID)
	})

	t.Run("TestExpiredJWT", func(t *testing.T) {
		claims.ExpiresAt = -expiration
		clientToken.Claims = claims

		ss, _ := clientToken.SignedString(secretBytes)

		verifier := KeyVerifier{ss, testDb}
		parsedKey, err = verifier.VerifyKey()

		assert.NotNil(t, err, "expired JWT should not verify")
	})

}
