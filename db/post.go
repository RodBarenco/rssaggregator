package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	FeedID      uuid.UUID      `gorm:"type:uuid"`
	Title       string         `gorm:"not null"`
	Description sql.NullString `gorm:"type:text"`
	Url         string         `gorm:"unique;not null"`
	PublishedAt time.Time      `gorm:"not null"`
}

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FeedID      uuid.UUID
	Title       string
	Description sql.NullString
	Url         string
	PublishedAt time.Time
}

func CreatePost(ctx context.Context, db *gorm.DB, arg CreatePostParams) (Post, error) {
	// Verificar se o arg.Url já existe no banco de dados
	var existingPost Post
	if err := db.WithContext(ctx).Where("url = ?", arg.Url).First(&existingPost).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return Post{}, err
		}
	} else {
		return existingPost, nil
	}

	post := Post{
		ID:          arg.ID,
		CreatedAt:   arg.CreatedAt,
		UpdatedAt:   arg.UpdatedAt,
		FeedID:      arg.FeedID,
		Title:       arg.Title,
		Description: arg.Description,
		Url:         arg.Url,
		PublishedAt: arg.PublishedAt,
	}

	if err := db.WithContext(ctx).Create(&post).Error; err != nil {
		return Post{}, err
	}

	return post, nil
}

func GetAllPosts(ctx context.Context, db *gorm.DB) ([]Post, error) {
	var posts []Post

	err := db.WithContext(ctx).
		Order("published_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func GetFilteredPosts(ctx context.Context, db *gorm.DB, userID uuid.UUID, limit int) ([]Post, error) {
	var posts []Post

	err := db.WithContext(ctx).
		Joins("JOIN feed_follows ON feed_follows.feed_id = posts.feed_id").
		Where("feed_follows.user_id = ?", userID).
		Order("posts.published_at DESC").
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}
