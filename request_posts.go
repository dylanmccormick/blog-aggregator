package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dylanmccormick/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {

	input := r.URL.Query().Get("limit")

	limit := 5
	var err error

	if input != "" {
		limit, err = strconv.Atoi(input)
		if err != nil {
			log.Printf("Error converting string to int: %v, using default limit", input)
			log.Printf("ERR: %v", err)
			limit = 5
		}
	}

	params := database.GetPostsByUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	}

	resp, err := cfg.DB.GetPostsByUser(context.TODO(), params)
	if err != nil {
		log.Printf("ERROR: unable to get posts for user: %v", user.ID)
		respondWithError(w, 500, "Internal server error")
		return
	}

	type Response struct {
		Name        string    `json:"name"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Url         string    `json:"url"`
		PublishedAt time.Time `json:"time"`
		ID          uuid.UUID `json:"id"`
	}

	response := []Response{}

	for _, post := range resp {
		desc, err := post.Description.Value()
		if err != nil {
			log.Printf("error converting description: %v\n", err)
		}
		str, ok := desc.(string)
		if !ok {
			log.Printf("incorrect type for description. Not of type string\n")
		}

		pub := post.PublishedAt.Time
		response = append(response, Response{
			Name:        post.Name,
			Title:       post.Title,
			Description: str,
			Url:         post.Url,
			PublishedAt: pub,
			ID:          post.ID,
		})

	}

	respondWithJSON(w, 200, response)

}
