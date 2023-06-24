package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedFollows struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	UserID    uuid.UUID `gorm:"type:uuid;foreignKey:User_id;constraint:OnDelete:CASCADE"`
	FeedID    uuid.UUID `gorm:"type:uuid;foreignKey:Feed_id;constraint:OnDelete:CASCADE;uniqueIndex:idx_user_feed_unique"`
}

type InsertFeedFollowsParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

func InsertFeedFollow(ctx context.Context, db *gorm.DB, arg InsertFeedFollowsParams) (FeedFollows, error) {
	feedfollows := FeedFollows{
		Model:  gorm.Model{},
		ID:     arg.ID,
		UserID: arg.UserID,
		FeedID: arg.FeedID,
	}

	if err := db.WithContext(ctx).Create(&feedfollows).Error; err != nil {
		return FeedFollows{}, err
	}

	return feedfollows, nil
}

func DeleteFeedFollow(ctx context.Context, db *gorm.DB, id uuid.UUID, userID uuid.UUID) (string, error) {
	feedfollows := FeedFollows{}
	if err := db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&feedfollows).Error; err != nil {
		return "", err
	}

	if err := db.WithContext(ctx).Delete(&feedfollows).Error; err != nil {
		return "", err
	}

	return "Feed_follow deleted", nil
}

func GetFeedFollows(ctx context.Context, db *gorm.DB, userID uuid.UUID) ([]FeedFollows, error) {
	var feedFollows []FeedFollows
	if err := db.WithContext(ctx).Where("user_id = ?", userID).Find(&feedFollows).Error; err != nil {
		return nil, err
	}
	return feedFollows, nil
}
