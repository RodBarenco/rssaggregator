package main

import (
	"time"

	"github.com/google/uuid"
)

// USER CREATION
type CreateUserResponse struct {
	User CreatedResponse `json:"user"`
}
type CreatedResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	ApiKey string    `json:"api_key"`
}

// GET BY APIKEY
type GetUserResponse struct {
	User UserGetedResponse `json:"user"`
}
type UserGetedResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// CREATE FEED
type CreateFeedResponse struct {
	User FeedCreatedResponse `json:"feed"`
}
type FeedCreatedResponse struct {
	Name   string    `json:"name"`
	Url    string    `json:"Url"`
	UserID uuid.UUID `json:"user_id"`
}

// GET FEEDS
type FeedResponse struct {
	Name   string    `json:"name"`
	URL    string    `json:"url"`
	UserID uuid.UUID `json:"user_id"`
}

// GET FEEDFOLLOWS
type FeedFollowResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	FeedIDs    uuid.UUID `json:"feed_ids"`
	FollowedAt time.Time `json:"updated_at"`
}

// GET POSTS FROM USERS

type GetFilteredPostsResponse struct {
	Posts []FilteredPost `json:"posts"`
}

type FilteredPost struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	FeedID      uuid.UUID `json:"feed_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Url         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
}
