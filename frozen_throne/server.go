package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/TheJokersThief/frozen-throne/frozen_throne/config"
)

type StatusResponse struct {
	Frozen bool  `json:"frozen"`
	Error  error `json:"error,omitempty"`
}

// IngestHTTP handles routing the HTTP request of the cloud function
func IngestHTTP(w http.ResponseWriter, r *http.Request) {
	ft := NewFrozenThrone(context.Background())

	switch r.Method {
	case http.MethodGet: // Status check
		query := r.URL.Query()
		token, tokenErr := getQueryParam(query, "token")
		if tokenErr != nil {
			http.Error(w, tokenErr.Error(), http.StatusBadRequest)
			return
		}

		if verifyWriteToken(token, ft.Config) || verifyReadOnlyToken(token, ft.Config) {
			// If at least one of the tokens is correct
			repo, repoErr := getQueryParam(query, "repo")
			if repoErr != nil {
				http.Error(w, repoErr.Error(), http.StatusBadRequest)
				return
			}
			_, statusErr := ft.Check(repo)
			if statusErr == nil {
				// If the status error is nil, that means it exists
				json.NewEncoder(w).Encode(StatusResponse{Frozen: true})
				return
			} else {
				json.NewEncoder(w).Encode(StatusResponse{Frozen: false, Error: statusErr})
				return
			}
		} else {
			http.Error(w, "401 - Unauthorized: Bad token", http.StatusBadRequest)
			return
		}
	case http.MethodPost: // Freeze
		token := r.FormValue("token")
		repo := r.FormValue("repo")
		user := r.FormValue("user")
		if token == "" || repo == "" || user == "" {
			http.Error(w, "400 - Bad Request: Missing token, repo or user", http.StatusBadRequest)
			return
		}

		if verifyWriteToken(token, ft.Config) {
			freezeErr := ft.Freeze(repo)
			if freezeErr != nil {
				json.NewEncoder(w).Encode(StatusResponse{Frozen: false, Error: freezeErr})
				return
			}
			json.NewEncoder(w).Encode(StatusResponse{Frozen: true})
			return
		} else {
			http.Error(w, "401 - Unauthorized: Bad token", http.StatusBadRequest)
			return
		}

	case http.MethodPatch: // Unfreeze
		token := r.FormValue("token")
		repo := r.FormValue("repo")
		user := r.FormValue("user")
		if token == "" || repo == "" || user == "" {
			http.Error(w, "400 - Bad Request: Missing token, repo or user", http.StatusBadRequest)
			return
		}

		if verifyWriteToken(token, ft.Config) {
			freezeErr := ft.Unfreeze(repo)
			if freezeErr != nil {
				json.NewEncoder(w).Encode(StatusResponse{Frozen: true, Error: freezeErr})
				return
			}
			json.NewEncoder(w).Encode(StatusResponse{Frozen: false})
			return
		} else {
			http.Error(w, "401 - Unauthorized: Bad token", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

// verifyWriteToken makes sure the token provided matches the token in the config
func verifyWriteToken(token string, config config.Config) bool {
	return token == config.WriteSecret
}

// verifyReadOnlyToken makes sure the token provided matches the token in the config
func verifyReadOnlyToken(token string, config config.Config) bool {
	return token == config.ReadOnlySecret
}

// getQueryParam extracts a single value for a given query parameter
func getQueryParam(query url.Values, key string) (string, error) {
	param, paramExists := query[key]
	if !paramExists {
		return "", fmt.Errorf("400 - Bad Request: ?%s GET parameter not supplied, bad request", key)
	}

	return param[0], nil
}

func main() {
	// handle route using handler function
	http.HandleFunc("/", IngestHTTP)

	// listen to port
	http.ListenAndServe(":5050", nil)
}
