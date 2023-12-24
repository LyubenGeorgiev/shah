package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/LyubenGeorgiev/shah/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func SetupDatabase() {
	// Fetch environment variables
	host := getenv("DATABASE_HOST", "localhost")
	portStr := getenv("POSTGRES_PORT", "5432")
	user := getenv("POSTGRES_USER", "root")
	password := getenv("POSTGRES_PASSWORD", "root")
	dbname := getenv("POSTGRES_DB", "postgres")

	// Convert port string to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Error converting port to integer:", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=EET", host, user, password, dbname, port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate will create tables based on provided structs
	err = db.AutoMigrate(&models.Game{}, &models.Tournament{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tables created successfully!")
}

func GetDB() *gorm.DB {
	return db
}
