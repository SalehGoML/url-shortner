package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SalehGoML/config"
	"github.com/SalehGoML/internal/handlers"
	"github.com/SalehGoML/internal/middleware"
	"github.com/SalehGoML/internal/repository"
	"github.com/SalehGoML/internal/service"

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

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	authHandler := handlers.NewAuthHandler(authService)

	mux := http.NewServeMux()

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	mux.Handle(
		"/dashboard",
		middleware.AuthMiddleware(
			cfg.JWTSecret,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("welcome"))
			}),
		),
	)

	log.Println("server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
