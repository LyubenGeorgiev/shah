package db

import "github.com/LyubenGeorgiev/shah/db/models"

type Storage interface {
	CreateUser(*models.User) error
	SaveUser(*models.User) error
	FindOneUser(email string, password string) (uint, error)
	FindByUserID(userID string) (*models.User, error)
	UpdateUserImage(userID string, image string) error
	FetchUsersByUsername(username string) ([]models.User, error)
	CreateGame(game *models.Game) error
	GetGame(gameID string) (*models.Game, error)
	CreateNews(*models.News) error
	GetAllNews() ([]models.News, error)
	GetMatchHistoryGames(userID string, page int, limit int) ([]models.Game, error)

	CreateMessage(msg *models.Message) error
	GetRecentChatsUserIDs(userID string, page int, limit int) ([]uint, error)
	GetRecentMessagesWith(userID1, userID2 string, page int, limit int) ([]models.Message, error)
	GetAllUsers(page, limit int) ([]models.User, error)
	DeleteUserByID(userID uint) error
	UpdateUser(userID uint, role string) error
}
