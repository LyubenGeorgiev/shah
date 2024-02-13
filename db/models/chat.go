package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model  `json:"-"`
	Use         uint    `json:"username" gorm:"uniqueIndex"`
	Email       string  `json:"email" gorm:"type:varchar(100);uniqueIndex"`
	Password    string  `json:"password"`
	Rating      float64 `json:"rating"`
	GamesPlayed int     `json:"gamesPlayed"`
	GamesWon    int     `json:"gamesWon"`
	Image       string  `json:"image"`
	Role        string  `json:"role"`
}
