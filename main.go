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

	return mux

}
