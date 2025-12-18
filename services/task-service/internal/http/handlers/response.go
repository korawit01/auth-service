package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

var (
	ErrMissingHeader = errors.New("missing X-User-ID")
	ErrUnauthorized  = errors.New("unauthorized")
)

func parseUserID(r *http.Request) (string, error) {
	id := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if id == "" {
		return "", ErrUnauthorized
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrUnauthorized) {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	log.Printf("handler error: %v", err)
	writeError(w, http.StatusInternalServerError, err.Error())
}
