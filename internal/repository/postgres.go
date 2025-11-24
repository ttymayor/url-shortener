package repository

import (
	"fmt"
	"log"

	"github.com/ttymayor/url-shortener/internal/config"
	"github.com/ttymayor/url-shortener/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresClient(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.Port,
		cfg.Postgres.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// Auto Migrate
	if err := db.AutoMigrate(&model.URL{}); err != nil {
		log.Println("failed to auto migrate database: " + err.Error())
	}

	return db
}
