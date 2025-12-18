package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var ErrMissingUserID = errors.New("missing X-User-ID")

func parseUserID(r *http.Request) (string, error) {
	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		return "", ErrMissingUserID
	}
	return userID, nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
