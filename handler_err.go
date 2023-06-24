package main

import "net/http"

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}

type gormError struct {
	Error string `json:"error"`
}

// MarshalJSON converte o objeto gormError para o formato JSON usado pelo GORM.
func (g gormError) MarshalJSON() ([]byte, error) {
	return []byte(`{"error": "` + g.Error + `"}`), nil
}
