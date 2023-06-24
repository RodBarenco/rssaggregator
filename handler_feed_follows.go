package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	database "github.com/RodBarenco/rssaggregator/db"
	"github.com/go-chi/chi"

	"github.com/google/uuid"
)

func (apiCfg *APIapiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err.Error()))
		return
	}

	feedfollows, err := database.InsertFeedFollow(r.Context(), apiCfg.DB, database.InsertFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feedfollow: %s", err.Error()))
		return
	}

	res := FeedFollowResponse{
		ID:         feedfollows.ID,
		UserID:     feedfollows.UserID,
		FeedIDs:    feedfollows.FeedID,
		FollowedAt: feedfollows.UpdatedAt,
	}

	respondWithJSON(w, 201, res)
}

func (apiCfg *APIapiConfig) handlerGetFeedsFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedfollows, err := database.GetFeedFollows(r.Context(), apiCfg.DB, user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed_follows: %s", err.Error()))
		return
	}
	res := make([]FeedFollowResponse, len(feedfollows))
	for i, feedfollows := range feedfollows {
		res[i] = FeedFollowResponse{
			ID:         feedfollows.ID,
			UserID:     feedfollows.UserID,
			FeedIDs:    feedfollows.FeedID,
			FollowedAt: feedfollows.UpdatedAt,
		}
	}
	respondWithJSON(w, 200, res)
}

func (apiCfg *APIapiConfig) handlerDeleteFeedsFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err.Error()))
		return
	}

	deleteResult, err := database.DeleteFeedFollow(r.Context(), apiCfg.DB, feedFollowID, user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed_follows: %s", err.Error()))
		return
	}

	res := deleteResult
	respondWithJSON(w, 200, res)
}
