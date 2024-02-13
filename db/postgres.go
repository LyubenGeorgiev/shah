package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/LyubenGeorgiev/shah/db/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func NewPostgresStorage() *PostgresStorage {
	// Fetch environment variables
	host := Getenv("DATABASE_HOST", "localhost")
	portStr := Getenv("POSTGRES_PORT", "5432")
	user := Getenv("POSTGRES_USER", "root")
	password := Getenv("POSTGRES_PASSWORD", "root")
	dbname := Getenv("POSTGRES_DB", "postgres")

	// Convert port string to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Error converting port to integer:", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=EET", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate will create tables based on provided structs
	err = db.AutoMigrate(&models.Game{}, &models.Tournament{}, &models.User{}, &models.News{}, &models.Message{})
	if err != nil {
		log.Fatal(err)
	}

	return &PostgresStorage{
		db: db,
	}
}

func (ps *PostgresStorage) CreateUser(user *models.User) error {
	return ps.db.Create(user).Error
}

func (ps *PostgresStorage) SaveUser(user *models.User) error {
	return ps.db.Save(user).Error
}

func (ps *PostgresStorage) FindOneUser(email, password string) (uint, error) {
	user := &models.User{}

	if err := ps.db.Where("email = ?", email).First(user).Error; err != nil {
		return 0, fmt.Errorf("Wrong email or password")
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return 0, fmt.Errorf("Wrong email or password")
	}

	return user.ID, nil
}

func (ps *PostgresStorage) FindByUserID(userID string) (*models.User, error) {
	var user models.User
	if err := ps.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserImage updates the image of the user with the given ID
func (ps *PostgresStorage) UpdateUserImage(userID string, image string) error {
	return ps.db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("image", image).Error
}

func (ps *PostgresStorage) FetchUsersByUsername(username string) ([]models.User, error) {
	var users []models.User
	if err := ps.db.Limit(5).Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (ps *PostgresStorage) CreateGame(game *models.Game) error {
	return ps.db.Create(game).Error
}

func (ps *PostgresStorage) GetGame(gameID string) (*models.Game, error) {
	var game models.Game
	if err := ps.db.First(&game, "id = ?", gameID).Error; err != nil {
		return nil, err
	}

	return &game, nil
}

func (ps *PostgresStorage) CreateNews(news *models.News) error {
	return ps.db.Create(news).Error
}

func (ps *PostgresStorage) GetAllNews() ([]models.News, error) {
	var newsList []models.News
	if err := ps.db.Find(&newsList).Error; err != nil {
		return nil, err
	}
	return newsList, nil
}

func (ps *PostgresStorage) GetMatchHistoryGames(userID string, page int, limit int) ([]models.Game, error) {
	var games []models.Game

	if err := ps.db.Limit(limit).Offset(page*limit).Where("white_id = ?", userID).Or("black_id = ?", userID).Order("created_at DESC").Find(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (ps *PostgresStorage) CreateMessage(msg *models.Message) error {
	return ps.db.Create(msg).Error
}

func (ps *PostgresStorage) GetRecentChatsUserIDs(userID string, page int, limit int) ([]uint, error) {
	var msgs []models.Message
	err := ps.db.Raw("select DISTINCT LEAST(\"from\", \"to\") as from, GREATEST(\"from\", \"to\") as to, created_at FROM messages where \"from\" = ? or \"to\" = ? order by created_at DESC limit ? offset ?", userID, userID, limit, page*limit).Scan(&msgs).Error
	if err != nil {
		return nil, err
	}

	res := []uint{}
	user, _ := strconv.ParseInt(userID, 10, 64)

	for _, msg := range msgs {
		if int64(msg.From) == user {
			res = append(res, msg.To)
		} else {
			res = append(res, msg.From)
		}
	}

	return res, nil
}

func (ps *PostgresStorage) GetRecentMessagesWith(userID1, userID2 string, page int, limit int) ([]models.Message, error) {
	var msgs []models.Message
	err := ps.db.Raw("select * FROM messages where (\"from\" = ? and \"to\" = ?) or (\"from\" = ? and \"to\" = ?) order by created_at DESC limit ? offset ?", userID1, userID2, userID2, userID1, limit, page*limit).Scan(&msgs).Error
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
