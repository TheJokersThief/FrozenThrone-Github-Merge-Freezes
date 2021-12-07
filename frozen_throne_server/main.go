package frozen_throne_server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TheJokersThief/frozen-throne/frozen_throne"
	"github.com/TheJokersThief/frozen-throne/frozen_throne/config"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type StatusResponse struct {
	Frozen bool  `json:"frozen"`
	Error  error `json:"error,omitempty"`
}

var ft *frozen_throne.FrozenThrone
var serverConfig *Config

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Access-Token")

		if verifyWriteToken(token, ft.Config) {
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func FreezeHandler(w http.ResponseWriter, r *http.Request) {
	repo := mux.Vars(r)["repo"]
	token := r.FormValue("token")
	user := r.FormValue("user")

	if token == "" || repo == "" || user == "" {
		http.Error(w, "400 - Bad Request: Missing token, repo or user", http.StatusBadRequest)
		return
	}

	freezeErr := ft.Freeze(repo)
	if freezeErr != nil {
		json.NewEncoder(w).Encode(StatusResponse{Frozen: false, Error: freezeErr})
		return
	}
	json.NewEncoder(w).Encode(StatusResponse{Frozen: true})
}

func UnfreezeHandler(w http.ResponseWriter, r *http.Request) {
	repo := mux.Vars(r)["repo"]
	token := r.FormValue("token")
	user := r.FormValue("user")

	if token == "" || repo == "" || user == "" {
		http.Error(w, "400 - Bad Request: Missing token, repo or user", http.StatusBadRequest)
		return
	}

	freezeErr := ft.Unfreeze(repo)
	if freezeErr != nil {
		json.NewEncoder(w).Encode(StatusResponse{Frozen: true, Error: freezeErr})
		return
	}
	json.NewEncoder(w).Encode(StatusResponse{Frozen: false})
}

// verifyWriteToken makes sure the token provided matches the token in the config
func verifyWriteToken(token string, config config.Config) bool {
	return token == config.WriteSecret
}

func Main() {
	ft = frozen_throne.NewFrozenThrone(context.Background())
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/{repo}/freeze", FreezeHandler).Subrouter().Use(AuthMiddleware)
	r.HandleFunc("/{repo}/unfreeze", UnfreezeHandler).Subrouter().Use(AuthMiddleware)
	r.HandleFunc("/{repo}/github-webhook", WebhookHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
