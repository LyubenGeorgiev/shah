package models

import (
	"gorm.io/gorm"
)

// User struct declaration
type User struct {
	gorm.Model  `json:"-"`
	Username    string  `json:"username" gorm:"uniqueIndex"`
	Email       string  `json:"email" gorm:"type:varchar(100);uniqueIndex"`
	Password    string  `json:"password"`
	Rating      float64 `json:"rating"`
	GamesPlayed int     `json:"gamesPlayed"`
	GamesWon    int     `json:"gamesWon"`
	Image       string  `json:"image"`
	Role 		string  `json:"role"`
}
