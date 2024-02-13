package models

import (
	"time"
)

type Message struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	From      uint      `gorm:"index"`
	To        uint      `gorm:"index"`
	Text      string    `json:"text"`
}
