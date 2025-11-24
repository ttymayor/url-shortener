package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ttymayor/url-shortener/internal/model"
	"gorm.io/gorm"
)

type URLRepository interface {
	Save(url *model.URL) error
	FindByShortCode(code string) (*model.URL, error)
	FindAll() ([]*model.URL, error)
}

type urlRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewURLRepository(db *gorm.DB, redis *redis.Client) URLRepository {
	return &urlRepository{
		db:    db,
		redis: redis,
	}
}

func (r *urlRepository) Save(url *model.URL) error {
	if err := r.db.Create(url).Error; err != nil {
		return err
	}

	// Cache the new URL
	ctx := context.Background()
	r.redis.Set(ctx, url.ShortCode, url.OriginalURL, 24*time.Hour)

	return nil
}

func (r *urlRepository) FindByShortCode(code string) (*model.URL, error) {
	ctx := context.Background()

	// Try Redis first
	val, err := r.redis.Get(ctx, code).Result()
	if err == nil {
		return &model.URL{
			ShortCode:   code,
			OriginalURL: val,
		}, nil
	}

	// If not in Redis, check DB
	var url model.URL
	if err := r.db.Where("short_code = ?", code).First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Cache it for next time
	r.redis.Set(ctx, url.ShortCode, url.OriginalURL, 24*time.Hour)

	return &url, nil
}

func (r *urlRepository) FindAll() ([]*model.URL, error) {
	var urls []*model.URL
	if err := r.db.Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}
