package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Name      string    `gorm:"not null"`
	ApiKey    string    `gorm:"not null"`
}

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func CreateUser(ctx context.Context, db *gorm.DB, arg CreateUserParams) (User, error) {
	apiKey, err := generateAPIKey()
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:     arg.ID,
		Name:   arg.Name,
		ApiKey: apiKey,
	}

	if err := db.WithContext(ctx).Create(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}
