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

func SetupDatabase() {
	// Fetch environment variables
	host := os.Getenv("DATABASE_HOST")
	portStr := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Convert port string to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Error converting port to integer:", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=EET search_path=Shah", host, user, password, dbname, port)

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
