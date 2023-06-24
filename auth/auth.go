package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIkey(headres http.Header) (string, error) {
	val := headres.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info foud")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part od auth header")
	}
	return vals[1], nil

}
