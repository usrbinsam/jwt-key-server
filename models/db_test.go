package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDb *gorm.DB

func setupTestDb() error {

	if testDb != nil {
		return nil
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return err
	}

	testDb = db
	return nil
}
