package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dylanmccormick/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareAuth(next authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("Authorization")

		if apiKey == "" {
			log.Printf("Unauthorized: no apikey provided")
			respondWithError(w, 401, "Unauthorized: no apikey provided")
			return
		}

		apiKey = strings.Split(apiKey, " ")[1]

		user, err := cfg.DB.GetUserByAPIKey(context.TODO(), apiKey)
		if err != nil {
			log.Printf("ERROR: unable to get user in db. apiKey: %v", apiKey)
			respondWithError(w, 500, "Internal Server Error")
			return
		}
		t := database.User{}

		if user == t {
			respondWithError(w, 404, "User not found")
			return
		}
		next(w, r, user)
	})
}
