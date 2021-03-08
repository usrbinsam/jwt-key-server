package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserModel(t *testing.T) {
	user := User{
		Username: "sheldon",
	}

	t.Run("TestShortPassword", func(t *testing.T) {
		assert.NotNil(t, user.SetPassword("foo"), "expected error when setting a short password")
	})

	t.Run("TestSetPassword", func(t *testing.T) {
		assert.Nil(t, user.SetPassword("password"), "error setting valid password")
	})

	t.Run("TestVerifyCorrectPassword", func(t *testing.T) {
		ok, err := user.VerifyPassword("password")
		assert.Nil(t, err, "VerifyPassword() failed")
		assert.True(t, ok, "incorrect password")
	})

	testDb.Save(&user)

	t.Run("TestUserLoaderBadPassword", func(t *testing.T) {
		_, err := UserLoader(testDb, "sheldon", "invalid")
		assert.NotNil(t, err, "expected error from invalid password")
	})

	t.Run("TestUserLoaderBadUsername", func(t *testing.T) {
		_, err := UserLoader(testDb, "sam", "password")
		assert.NotNil(t, err, "expected error from invalid username")
	})

	t.Run("TestUserLoaderValidCredentials", func(t *testing.T) {
		user, err := UserLoader(testDb, "sheldon", "password")
		assert.Nil(t, err)
		assert.NotNil(t, user, "expected non-nil user pointer")
		assert.Equal(t, user.Username, "sheldon")
	})
}
