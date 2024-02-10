package db

import "github.com/LyubenGeorgiev/shah/db/models"

type Storage interface {
	CreateUser(*models.User) error
	FindOneUser(email string, password string) (uint, error)
	FindByUserID(userID string) (*models.User, error)
	UpdateUserImage(userID string, image string) error
}
