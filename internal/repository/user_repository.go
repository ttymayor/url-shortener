package repository

import (
	"errors"

	"github.com/ttymayor/url-shortener/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(username, password string) error
	GetUserByUsername(username string) (*model.User, error)
	SeedAdminUser(username, password string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	// Auto migrate User model
	if err := db.AutoMigrate(&model.User{}); err != nil {
		panic("failed to migrate user model: " + err.Error())
	}
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// SeedAdminUser checks if the admin user exists.
// If not, it creates it.
// If it exists, it updates the password to match the config.
func (r *userRepository) SeedAdminUser(username, password string) error {
	existing, err := r.GetUserByUsername(username)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if existing == nil {
		// Create new admin
		user := &model.User{
			Username: username,
			Password: string(hashedPassword),
		}
		return r.db.Create(user).Error
	}

	// Update existing admin password (to ensure .env and DB are in sync)
	existing.Password = string(hashedPassword)
	return r.db.Save(existing).Error
}
