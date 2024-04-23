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

func (cfg *apiConfig) createRSSFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type Input struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	input := Input{}

	err := decoder.Decode(&input)
	if err != nil {
		log.Printf("ERROR: unable to decode JSON in createRSSFeed")
		respondWithError(w, 500, "Internal Server Error")
		return
	}

	feed := database.CreateRssFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      input.Name,
		Url:       input.Url,
		UserID:    user.ID,
	}

	resp, err := cfg.DB.CreateRssFeed(context.TODO(), feed)
	if err != nil {
		log.Printf("ERROR: %v", err)
		respondWithError(w, 400, "Url for given feed already exists")
		return
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    resp.ID,
	}

	feedFollowResp, err := cfg.DB.CreateFeedFollow(context.TODO(), feedFollow)

	if err != nil {
		log.Printf("ERROR: %v", err)
		respondWithError(w, 400, "Url for given feed already exists")
		return
	}

	type Output struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}
	output := Output{
		resp,
		feedFollowResp,
	}

	respondWithJSON(w, 200, output)
}

func (cfg *apiConfig) getRSSFeeds(w http.ResponseWriter, r *http.Request) {
	resp, err := cfg.DB.GetAllRssFeeds(context.TODO())

	if err != nil {
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 200, resp)
}
