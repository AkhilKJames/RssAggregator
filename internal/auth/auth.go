package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts api key from header
// Example :
// Authorization: ApiKey {api_key}
func GetAPIKey(h http.Header) (string, error) {
	val := h.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 || vals[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}
	return vals[1], nil
}
