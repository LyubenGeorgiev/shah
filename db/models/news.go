package models

import (
	"gorm.io/gorm"
)

// User struct declaration
type News struct {
	gorm.Model
	Image       string
	Title       string
	Description string
	URL         string
}
