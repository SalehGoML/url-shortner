package repository

import (
	"github.com/SalehGoML/internal/models"
	"gorm.io/gorm"
)

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

func (r *urlRepository) GetByShortCode(code string) (*models.URL, error) {
	var url models.URL
	err := r.db.
		Where("short_code = ? AND is_active = ?", code, true).
		First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) ListByUser(userID uint) ([]models.URL, error) {
	var urls []models.URL
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&urls).Error
	return urls, err
}

func (r *urlRepository) IncrementClicks(id uint) error {
	return r.db.Model(&models.URL{}).
		Where("id = ?", id).
		UpdateColumn("clicks", gorm.Expr("clicks + 1")).
		Error
}

func (r *urlRepository) Deactivate(id uint) error {
	return r.db.Model(&models.URL{}).
		Where("id = ?", id).
		Update("is_active", false).
		Error
}
