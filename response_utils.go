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
