package model

import (
    "fmt"
    "log"
    "gorm.io/driver/postgres" 
    "gorm.io/gorm"
)

var db *gorm.DB

func SetupDatabase() {
    
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5433 sslmode=disable TimeZone=EET search_path=Shah"

    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // AutoMigrate will create tables based on provided structs
    err = db.AutoMigrate(&Game{}, &Tournament{}) =
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Tables created successfully!")
}

func GetDB() *gorm.DB {
    return db
}