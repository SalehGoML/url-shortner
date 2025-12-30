package service

import (
	"github.com/SalehGoML/internal/models"
	"github.com/SalehGoML/internal/repository"
)

type AuthService interface {
	Register(name, email, password string) (*models.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}
