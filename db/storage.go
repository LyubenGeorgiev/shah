package db

import "github.com/LyubenGeorgiev/shah/db/models"

type Storage interface {
	CreateUser(*models.User) error
	FindOneUser(email string, password string) (uint, error)
}
