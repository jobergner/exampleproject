package server

import (
	"encoding/json"
	"exampleproject/api"
	"exampleproject/log"
	"io"
	"net/http"
)

func newHandler() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/create-quiz", createQuizHandler)
	handler.HandleFunc("/login", login)
	handler.HandleFunc("/logout", logout)

	handler.Handle("/secret", authMiddleware(http.HandlerFunc(secret)))
	return handler
}

func createQuizHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Log(log.ReadBody)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	var quiz api.NewQuiz
	if err := json.Unmarshal(body, &quiz); err != nil {
		log.Log(log.Unmarshal, log.JSONData(body))
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	if err := api.CreateQuiz(r.Context(), quiz); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
