package serve

import (
	"encoding/json"
	"exampleproject/action"
	"exampleproject/errmsg"
	"exampleproject/log"
	"io"
	"net/http"
)

func newHandler() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/create-quiz", createQuizHandler)
	return handler
}

func createQuizHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Log(errmsg.ReadBody("createQuizHandler"))
		return
	}

	var quiz action.NewQuiz
	if err := json.Unmarshal(body, &quiz); err != nil {
		log.Log(errmsg.Unmarshal("NewQuiz"), log.JSONData(body))
		return
	}

	if err := action.CreateQuiz(r.Context(), quiz); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
}
