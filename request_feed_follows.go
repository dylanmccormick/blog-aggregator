package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dylanmccormick/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) followFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type Input struct {
		FeedId string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	input := Input{}

	err := decoder.Decode(&input)
	if err != nil {
		log.Printf("ERROR: unable to decode JSON in followFeed")
		respondWithError(w, 500, "Internal Server Error")
		return
	}

	if input.FeedId == "" {
		respondWithError(w, 400, "No valid input given")
		return
	}

	feedId, err := uuid.Parse(input.FeedId)
	if err != nil {
		log.Printf("ERROR: unable to parse UUID in followFeed")
		respondWithError(w, 500, "Internal Server Error")
		return
	}

	cff := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedId,
	}

	resp, err := cfg.DB.CreateFeedFollow(context.TODO(), cff)
	if err != nil {
		log.Printf("Error: unable to create feed follow: %v", err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 200, resp)
}

func (cfg *apiConfig) unfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.PathValue("feedFollowId") == "" {
		respondWithError(w, 400, fmt.Sprintf("Expected an id value, recieved: %v", r.PathValue("id")))

	}

	input := r.PathValue("feedFollowId")

	if input == "" {
		respondWithError(w, 400, "No valid input given")
		return
	}

	feedId, err := uuid.Parse(input)
	if err != nil {
		log.Printf("ERROR: unable to parse UUID in followFeed")
		respondWithError(w, 500, "Internal Server Error")
		return
	}

	err = cfg.DB.DeleteFeedFollow(context.TODO(), feedId)
	if err != nil {
		log.Printf("ERROR: unalbe to delete feed: %v", input)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 200, "")

}
func (cfg *apiConfig) getFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	resp, err := cfg.DB.GetFeedFollows(context.TODO(), user.ID)
	if err != nil {
		log.Printf("ERROR: unable to get follows for user: %v", user.ID)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 200, resp)

}
