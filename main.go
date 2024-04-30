package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/dylanmccormick/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("POSTGRES_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to open postgres db")
	}

	dbQueries := database.New(db)
	cfg := apiConfig{
		dbQueries,
	}
	log.Printf("Running feed worker\n")
	go cfg.feedWorker(3)

	mux := cfg.mux()

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Handler: corsMux,
		Addr:    ":" + port,
	}

	log.Printf("Starting server on port %s...", port)

	log.Fatal(server.ListenAndServe())

}

func (cfg *apiConfig) mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/readiness", requestReadiness)
	mux.HandleFunc("GET /v1/error", requestError)
	mux.HandleFunc("POST /v1/users", cfg.createUser)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.getUsers))
	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.createRSSFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.getRSSFeeds)
	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.followFeed))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowId}", cfg.middlewareAuth(cfg.unfollowFeed))
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.getFeedFollows))
	mux.HandleFunc("GET /v1/posts", cfg.middlewareAuth(cfg.getPostsByUser))

	return mux

}
