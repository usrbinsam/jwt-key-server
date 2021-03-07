package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"index:uniqueIndex"` // Unique login name.
	Password []byte `gorm:"type:char(80)"`     // Not meant to be access directly. Use SetPassword() and VerifyPassword().
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

	user.Password = bytes
	return nil
}

// VerifyPassword compares a given plaintext password to the hashed
// password on the user.
func (user *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password)) == nil
}

// Save acts as an upsert, calling Create if the ID property is 0, otherwise Save.
func (user *User) Save(db *gorm.DB) {
	if user.ID == 0 {
		db.Create(user)
	} else {
		db.Save(user)
	}

	db.Commit()
}

// UserLoader acts as a login function. Returns a user given a username and password.
func UserLoader(db *gorm.DB, username, password string) (*User, error) {
	var user User

	if db.First(&user, "username = ?", username).Error != nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.VerifyPassword(password) {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}
