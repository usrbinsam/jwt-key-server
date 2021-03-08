package models

import (
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"index:uniqueIndex"` // Unique login name.
	Password string // Not meant to be access directly. Use SetPassword() and VerifyPassword().
}

// SetPassword hashes a given password using the bcrypt extension. Returns any error returned from
// bcrypt.GenerateFromPassword(), or if the password does not meet minimum complexity requirements.
func (user *User) SetPassword(password string) error {
	if len(password) < 8 {
		return errors.New("insufficient password length")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return err
	}

	user.Password = hex.EncodeToString(bytes)
	return nil
}

// VerifyPassword compares a given plaintext password to the hashed
// password on the user.
func (user *User) VerifyPassword(password string) (bool, error) {
	b, err := hex.DecodeString(user.Password)
	if err != nil {
		return false, err
	}

	return bcrypt.CompareHashAndPassword(b, []byte(password)) == nil, nil
}

// UserLoader acts as a login function. Returns a user given a username and password.
func UserLoader(db *gorm.DB, username, password string) (*User, error) {
	var user User

	if db.First(&user, "username = ?", username).Error != nil {
		return nil, errors.New("invalid username or password")
	}

	ok, err := user.VerifyPassword(password)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}
