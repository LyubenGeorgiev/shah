package models

import "time"

type Chat struct {
	From      uint      `gorm:"primarykey"`
	To        uint      `gorm:"primaryKey"`
	UpdatedAt time.Time `gorm:"index"`
}
