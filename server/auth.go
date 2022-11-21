package server

import (
	"encoding/json"
	"exampleproject/api"
	"exampleproject/log"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key            = []byte("super-secret-key")
	store          = sessions.NewCookieStore(key)
	authCookieName = "exampleproject_authenticated"
)

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, authCookieName)

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, authCookieName)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Log(log.ReadBody)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	var loginData api.LoginData
	if err := json.Unmarshal(body, &loginData); err != nil {
		log.Log(log.Unmarshal, log.JSONData(body))
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	authenticated, err := api.AuthenticateUser(r.Context(), loginData)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	if !authenticated {
		http.Error(w, "incorrect login data", http.StatusUnauthorized)
		return
	}

	session.Values["authenticated"] = true
	store.Save(r, w, session)
	w.WriteHeader(http.StatusOK)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, authCookieName)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = false
	store.Save(r, w, session)
}
