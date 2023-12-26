package db

import "github.com/LyubenGeorgiev/shah/db/models"

type Storage interface {
	CreateUser(*models.User) error
}
