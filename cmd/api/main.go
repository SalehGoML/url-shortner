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
	urlRepo := repository.NewURLRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	urlService := service.NewURLService(urlRepo)

	authHandler := handlers.NewAuthHandler(authService)
	urlHandler := handlers.NewURLHandler(urlService)

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

	//mux.HandleFunc("/shorten", urlHandler.Shorten)

	mux.Handle(
		"/shorten",
		middleware.AuthMiddleware(
			cfg.JWTSecret,
			http.HandlerFunc(urlHandler.Shorten),
		),
	)

	mux.Handle(
		"/my/urls",
		middleware.AuthMiddleware(
			cfg.JWTSecret,
			http.HandlerFunc(urlHandler.ListMyURLs),
		),
	)

	mux.HandleFunc("/", urlHandler.Redirect)

	log.Println("server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
