package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SalehGoML/config"
	"github.com/SalehGoML/internal/handlers"
	"github.com/SalehGoML/internal/middleware"
	"github.com/SalehGoML/internal/repository"
	"github.com/SalehGoML/internal/service"

	"github.com/SalehGoML/internal/models"

	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	router.Handle("/shorten", middleware.AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(urlHandler.Shorten))).Methods("POST")
	router.Handle("/my/urls", middleware.AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(urlHandler.ListMyURLs))).Methods("GET")
	router.Handle("/url/deactivate", middleware.AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(urlHandler.Deactivate))).Methods("POST")

	router.Handle("/dashboard", middleware.AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})))

	router.HandleFunc("/{code}", urlHandler.Redirect).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + cfg.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server running on port", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
