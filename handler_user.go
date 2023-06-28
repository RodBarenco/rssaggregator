package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	database "github.com/RodBarenco/rssaggregator/db"

	"github.com/google/uuid"
)

func (apiCfg *APIapiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
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

	user, err := database.CreateUser(r.Context(), apiCfg.DB, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %s", err.Error()))
		return
	}

	res := CreateUserResponse{
		User: CreatedResponse{
			ID:     user.ID,
			Name:   user.Name,
			ApiKey: user.ApiKey,
		},
	}

	respondWithJSON(w, 201, res)
}

func (apiCfg *APIapiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {

	res := GetUserResponse{
		User: UserGetedResponse{
			ID:   user.ID,
			Name: user.Name,
		},
	}
	respondWithJSON(w, 200, res)
}

func (apiCfg *APIapiConfig) handlerGetPostsFromUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := database.GetAllPosts(r.Context(), apiCfg.DB)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get posts: %v", err))
		return
	}

	// Cria a resposta
	res := GetFilteredPostsResponse{
		Posts: make([]FilteredPost, len(posts)),
	}

	for i, post := range posts {
		res.Posts[i] = FilteredPost{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			FeedID:      post.FeedID,
			Title:       post.Title,
			Description: nil, // Define inicialmente como nil
			Url:         post.Url,
			PublishedAt: post.PublishedAt,
		}

		// Define o campo Description apenas se for válido
		if post.Description.Valid {
			res.Posts[i].Description = &post.Description.String
		}
	}

	respondWithJSON(w, 200, res)
}
