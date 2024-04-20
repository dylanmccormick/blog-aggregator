package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Sending response with status code %d\n", status)
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding JSON: %v", payload)
		respondWithError(w, 500, "internal server error")
	}
	w.WriteHeader(status)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorBody struct {
		Error string `json:"error"`
	}
	resp := ErrorBody{msg}
	respondWithJSON(w, code, resp)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			log.Printf("Writing to header in middlewareCors function 1")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
