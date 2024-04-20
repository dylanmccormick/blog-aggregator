package main

import "net/http"

func requestReadiness(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string `json:"status"`
	}

	okResp := Response{"ok"}

	respondWithJSON(w, 200, okResp)

}
