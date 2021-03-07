package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKey(t *testing.T) {
	k := Key{}

	t.Run("EmptySecretBytes", func(t *testing.T) {
		b, err := k.GetSecretBytes()
		assert.Nil(t, err, "GetSecretBytes had error while empty")
		assert.Empty(t, b, "GetSecretBytes() should return empty byte slice with no secret set")
	})

	t.Run("SetRandomSecret", func(t *testing.T) {
		n, err := k.SetRandomSecret(32)
		assert.Nil(t, err, "failed to write key secret")
		assert.Equal(t, n, 32, "wrote 32 byte secret")
	})
}
