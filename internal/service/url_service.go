package service

import (
	"math/rand"
	"time"

	"github.com/ttymayor/url-shortener/internal/model"
	"github.com/ttymayor/url-shortener/internal/repository"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLService interface {
	ShortenURL(originalURL string) (*model.URL, error)
	GetOriginalURL(shortCode string) (string, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) ShortenURL(originalURL string) (*model.URL, error) {
	// Generate a random short code
	// In a real app, you'd check for collisions
	shortCode := generateShortCode(6)

	url := &model.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}

	if err := s.repo.Save(url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) GetOriginalURL(shortCode string) (string, error) {
	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}
	if url == nil {
		return "", nil
	}
	return url.OriginalURL, nil
}

func generateShortCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
