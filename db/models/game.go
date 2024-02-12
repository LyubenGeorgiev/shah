package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	_ "gorm.io/driver/postgres"
)

type Game struct {
	ID        string         `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time      `gorm:"index"`
	WhiteID   string         `gorm:"type:text"`
	BlackID   string         `gorm:"type:text"`
	WinnerID  sql.NullString `gorm:"type:text"`
	Moves     pq.Int32Array  `gorm:"type:integer[]"`
}
