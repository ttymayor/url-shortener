package service

import (
	"errors"
	"math/rand"
	"net/url"
	"regexp"
	"time"

	"github.com/ttymayor/url-shortener/internal/model"
	"github.com/ttymayor/url-shortener/internal/repository"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	ErrCodeInUse   = errors.New("short code already in use")
	ErrInvalidURL  = errors.New("invalid URL format")
	ErrInvalidCode = errors.New("invalid short code format or reserved word")
)

// Reserved words that cannot be used as short codes
var reservedCodes = map[string]bool{
	"api":    true,
	"health": true,
	"login":  true,
	"logout": true,
	"static": true,
	"admin":  true,
}

var codeRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

type URLService interface {
	ShortenURL(originalURL, customCode string) (*model.URL, error)
	GetOriginalURL(shortCode string) (string, error)
	GetAllURLs() ([]*model.URL, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) ShortenURL(originalURL, customCode string) (*model.URL, error) {
	// 1. Validate Original URL
	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, ErrInvalidURL
	}

	var shortCode string

	// 2. Handle Custom Code
	if customCode != "" {
		// Validate format (alphanumeric + hyphen)
		if !codeRegex.MatchString(customCode) {
			return nil, ErrInvalidCode
		}
		// Check reserved words
		if reservedCodes[customCode] {
			return nil, ErrInvalidCode
		}

		// Check if it already exists
		existing, err := s.repo.FindByShortCode(customCode)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrCodeInUse
		}
		shortCode = customCode
	} else {
		// Generate a random short code
		shortCode = generateShortCode(6)
		// In a real production app, you should also loop here to ensure the
		// generated code doesn't collide, though the probability is low.
	}

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

func (s *urlService) GetAllURLs() ([]*model.URL, error) {
	return s.repo.FindAll()
}

func generateShortCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
