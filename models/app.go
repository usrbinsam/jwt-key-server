package models

import (
	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	Keys           []Key
	Name           string
	SupportMessage string
}
