package main

import (
	"fmt"
	"net/http"

	"github.com/RodBarenco/rssaggregator/auth"
	database "github.com/RodBarenco/rssaggregator/db"
)

type autheHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *APIapiConfig) middlewareAuth(handler autheHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIkey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		var user database.User
		err = apiCfg.DB.Where("api_key = ?", apiKey).First(&user).Error
		if err != nil {
			respondWithError(w, 404, "User not found")
			return
		}

		handler(w, r, user)
	}
}
