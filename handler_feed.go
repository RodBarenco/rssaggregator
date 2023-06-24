package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	database "github.com/RodBarenco/rssaggregator/db"

	"github.com/google/uuid"
)

func (apiCfg *APIapiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err.Error()))
		return
	}

	if params.Name == "" {
		respondWithError(w, 400, "Invalid name")
		return
	}

	if params.URL == "" {
		respondWithError(w, 400, "Invalid url")
		return
	}

	feed, err := database.CreateFeed(r.Context(), apiCfg.DB, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		URL:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed: %s", err.Error()))
		return
	}

	res := CreateFeedResponse{
		User: FeedCreatedResponse{
			Name:   feed.Name,
			Url:    feed.Url,
			UserID: feed.UserID,
		},
	}

	respondWithJSON(w, 201, res)
}

func (apiCfg *APIapiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := database.GetFeeds(apiCfg.DB)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %s", err.Error()))
		return
	}

	res := make([]FeedResponse, len(feeds))
	for i, feed := range feeds {
		res[i] = FeedResponse{
			Name:   feed.Name,
			URL:    feed.Url,
			UserID: feed.UserID,
		}
	}

	respondWithJSON(w, 200, res)
}
