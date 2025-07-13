package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.HelloWorldHandler)

	r.Get("/auth/{provider}/callback", s.authCallbackHandler)
	r.Get("/auth/{provider}", s.authHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"message": "Hello World",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("Error handling auth callback: %v", err)
		http.Error(w, "Unable to log in", http.StatusUnauthorized)
		return
	}

	clientUrl := os.Getenv("CLIENT_URL")

	fmt.Println(user) // TODO: save to database

	http.Redirect(w, r, clientUrl, http.StatusFound)
}

func (s *Server) authHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	}

	gothic.BeginAuthHandler(w, r)
}
