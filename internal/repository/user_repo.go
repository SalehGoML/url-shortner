package repository

import "github.com/SalehGoML/internal/models"

type URLRepository interface {
	Create(url *models.URL) error
	GetByShortCode(shortCode string) (*models.URL, error)
	ListByUser(userID uint) ([]*models.URL, error)
	IncrementClicks(id uint) error
	Deactivate(id uint) error
}
