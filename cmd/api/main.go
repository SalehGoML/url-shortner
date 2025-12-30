package main

import (
	"fmt"
	"log"

	"github.com/SalehGoML/config"
	"github.com/SalehGoML/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	cfg := config.LoadConfig()
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	fmt.Println("Connected to DB")

	err = db.AutoMigrate(
		&models.User{},
		&models.URL{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("Tables migrated successfully")
}
