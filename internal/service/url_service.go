package service

import (
	"errors"
	"time"

	"github.com/SalehGoML/internal/models"
	"github.com/SalehGoML/internal/repository"
	"github.com/SalehGoML/internal/utils"
)

type URLService interface {
	Shorten(longURL string, userID uint, expiry *time.Time) (*models.URL, error)
	GetByCode(code string) (*models.URL, error)
	IncrementClicks(urlID uint) error
	ListByUser(userID uint) ([]*models.URL, error)
	Deactivate(urlID uint) error
}

type urlService struct {
	urlRepo repository.URLRepository
}

func NewURLService(urlRepo repository.URLRepository) URLService {
	return &urlService{urlRepo: urlRepo}
}

func (s *urlService) Shorten(
	longURL string,
	userID uint,
	expiry *time.Time,
) (*models.URL, error) {

	existingURLs, err := s.urlRepo.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	for _, u := range existingURLs {
		if u.LongURL == longURL && u.IsActive {
			return u, nil
		}
	}
	//fmt.Println("Shorten handler called")
	//

	var shortCode string
	for {
		shortCode = utils.GenerateShortCode(6)
		_, err := s.urlRepo.GetByShortCode(shortCode)
		if err != nil {
			break
		}
	}
	url := &models.URL{
		LongURL:   longURL,
		ShortCode: shortCode,
		UserID:    userID,
		Expiry:    expiry,
		IsActive:  true,
	}

	if err := s.urlRepo.Create(url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) GetByCode(code string) (*models.URL, error) {
	url, err := s.urlRepo.GetByShortCode(code)
	if err != nil {
		return nil, err
	}

	if url.Expiry != nil && url.Expiry.Before(time.Now()) {
		return nil, errors.New("link expired")
	}

	return url, nil
}

func (s *urlService) IncrementClicks(urlID uint) error {
	return s.urlRepo.IncrementClicks(urlID)
}

func (s *urlService) ListByUser(userID uint) ([]*models.URL, error) {
	return s.urlRepo.ListByUser(userID)
}

func (s *urlService) Deactivate(urlID uint) error {
	return s.urlRepo.Deactivate(urlID)
}
