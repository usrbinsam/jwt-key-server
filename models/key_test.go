package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	_ = setupTestDb()
	_ = testDb.AutoMigrate(Key{}, KeyActivation{})
}

func TestKey(t *testing.T) {
	k := Key{
		Status:         KeyStatusActive,
		MaxActivations: 1,
	}
	k.ID = 42

	var keySize = 64

	t.Run("EmptySecretBytes", func(t *testing.T) {
		b, err := k.GetSecretBytes()
		assert.Nil(t, err, "GetSecretBytes had error while empty")
		assert.Empty(t, b, "GetSecretBytes() should return empty byte slice with no secret set")
	})

	t.Run("SetRandomSecret", func(t *testing.T) {
		n, err := k.SetRandomSecret(keySize)
		assert.Nil(t, err, "failed to write key secret")
		assert.Equal(t, n, keySize, "wrote 32 byte secret")
	})

	t.Run("TestKeyActivation", func(t *testing.T) {
		machineName := "BRIDGES"

		ka := KeyActivation{
			Identifier: &machineName,
		}

		err := k.Activate(&ka)

		if assert.Nil(t, err, "Activate() failed") {
			assert.Equal(t, k.ID, ka.KeyID, "Activate() should have set the KeyActivation.KeyID field.")
		}

		k.MaxActivations = 2
		assert.NotNil(t, k.Activate(&ka), "Activate() should've failed because the identifier is already active")
	})

	t.Run("TestKeyMaxActivations", func(t *testing.T) {
		machineName := "FLOYD"
		ka := KeyActivation{
			Identifier: &machineName,
		}
		k.MaxActivations = 1
		err := k.Activate(&ka)

		assert.NotNil(t, err, "Activate() should've failed because it would exceed the max activations")
	})
}
