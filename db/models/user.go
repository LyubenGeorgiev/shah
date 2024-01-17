package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// User struct declaration
type User struct {
	gorm.Model `json:"-"`
	Username   string `json:"username"`
	Email      string `json:"email" gorm:"type:varchar(100);unique_index"`
	Password   string `json:"password"`
	Rating     float64 `json:"rating"`
	GamesPlayed int    `json:"gamesPlayed"`
	GamesWon    int    `json:"gamesWon"`
}


type Token struct {
	UserID         uint
	Username       string
	Email          string
	StandardClaims *jwt.StandardClaims
}
