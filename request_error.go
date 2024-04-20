package main

import "net/http"

func requestError(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error string `json:"error"`
	}

	errorResp := Response{"Internal Server Error"}

	respondWithJSON(w, 500, errorResp)

}
