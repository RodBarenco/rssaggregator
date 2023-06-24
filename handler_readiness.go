package main

import (
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, gormJSON{})
}

type gormJSON struct{}

// MarshalJSON converte o objeto gormJSON para o formato JSON usado pelo GORM.
func (g gormJSON) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}
