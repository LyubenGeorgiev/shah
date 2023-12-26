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
	err = db.AutoMigrate(&models.Game{}, &models.Tournament{}, &models.User{})
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

func (ps *PostgresStorage) FindOneUser(email, password string) error {
	user := &models.User{}

	if err := ps.db.Where("email = ?", email).First(user).Error; err != nil {
		return fmt.Errorf("Wrong email or password")
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return fmt.Errorf("Wrong email or password")
	}

	return nil
}
