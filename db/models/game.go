package models

import "time"

type Game struct {
	GameID          uint `gorm:"primaryKey"`
	GameName        string
	GameDescription string
	GameType        string
	GameCreatorID   uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Game) TableName() string {
	return "games"
}
