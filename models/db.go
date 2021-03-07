package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func setupDb() {
	dsn := os.Getenv("DSN")

	var err error

	db, err = gorm.Open(
		postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("setupDb() failed: " + err.Error())
	}
}
