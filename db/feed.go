package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feed struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	Name          string    `gorm:"not null"`
	Url           string    `gorm:"unique;not null"`
	UserID        uuid.UUID `gorm:"type:uuid;foreignKey:User_id;constraint:OnDelete:CASCADE"`
	LastFetchedAT time.Time
}

type CreateFeedParams struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	URL          string
	UserID       uuid.UUID
	LastFechedAT time.Time
}

func CreateFeed(ctx context.Context, db *gorm.DB, arg CreateFeedParams) (Feed, error) {
	feed := Feed{
		ID:     arg.ID,
		Name:   arg.Name,
		Url:    arg.URL,
		UserID: arg.UserID,
	}

	if err := db.WithContext(ctx).Create(&feed).Error; err != nil {
		return Feed{}, err
	}

	return feed, nil
}

func GetFeeds(db *gorm.DB) ([]Feed, error) {
	var feeds []Feed
	if err := db.Find(&feeds).Error; err != nil {
		return nil, err
	}
	return feeds, nil
}

func GetNextFeedsToFetch(ctx context.Context, limit int, db *gorm.DB) ([]Feed, error) {
	var feeds []Feed
	if err := db.Order("last_fetched_at ASC NULLS FIRST").Limit(limit).Find(&feeds).Error; err != nil {
		return nil, err
	}
	return feeds, nil
}

func MarkFeedFetched(ctx context.Context, db *gorm.DB, feedID uuid.UUID) (Feed, error) {
	var feed Feed
	if err := db.Model(&feed).Where("id = ?", feedID).Updates(map[string]interface{}{
		"last_fetched_at": time.Now(),
		"updated_at":      time.Now(),
	}).First(&feed).Error; err != nil {
		return Feed{}, err
	}
	return feed, nil
}
