package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dylanmccormick/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	input := Input{}

	err := decoder.Decode(&input)
	if err != nil {
		log.Printf("ERROR: unable to decode JSON in createUser")
		respondWithError(w, 500, "Internal Server Error")
		return
	}

	createUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      input.Name,
	}

	user, err := cfg.DB.CreateUser(
		context.TODO(),
		createUser,
	)
	if err != nil {
		log.Printf("ERROR: unable to create user in db. user: %v", createUser)
		respondWithError(w, 500, "Internal Server Error")
		return
	}

	respondWithJSON(w, 200, user)

}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, user)
}
