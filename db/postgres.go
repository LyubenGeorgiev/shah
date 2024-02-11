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
	err = db.AutoMigrate(&models.Game{}, &models.Tournament{}, &models.User{}, &models.News{})
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
	user := &models.User{}

	if err := ps.db.Where("id = ?", userID).First(user).Error; err != nil {
		return nil, fmt.Errorf("Wrong email or password")
	}

	return user, nil
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
