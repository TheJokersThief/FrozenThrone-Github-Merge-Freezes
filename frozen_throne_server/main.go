package frozen_throne_server

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
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

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Nothing here")
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
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	log.Println("Starting up")
	ft = frozen_throne.NewFrozenThrone(context.Background())
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", DefaultHandler)
	r.HandleFunc("/{repo}/freeze", FreezeHandler).Subrouter().Use(AuthMiddleware)
	r.HandleFunc("/{repo}/unfreeze", UnfreezeHandler).Subrouter().Use(AuthMiddleware)
	r.HandleFunc("/{repo}/github-webhook", WebhookHandler)

	log.Println("Handlers initialised")

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Beginning server")
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("Server running")

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	log.Println("Wait until shutdown signal")
	// Block until we receive our signal.
	<-c
	log.Println("Shutdown signal received")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Shutting down")
	os.Exit(0)
}
