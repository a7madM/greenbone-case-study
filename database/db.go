package database

import (
	"greenbone-case-study/models"
	"log"

	_ "github.com/glebarez/go-sqlite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	database, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "sqlite",
		DSN:        "computers.db",
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database ", err)
	}

	database.AutoMigrate(&models.Computer{})
	DB = database
}
