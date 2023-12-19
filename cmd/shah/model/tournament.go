package model

import "time"

type Tournament struct {
    TournamentID          uint   `gorm:"primaryKey"`
    TournamentName        string
    TournamentDescription string
    TournamentDate        time.Time
    OrganizerID           uint
    CreatedAt             time.Time
    UpdatedAt             time.Time
}


func (Tournament) TableName() string {
    return "Shah.tournaments" 
}