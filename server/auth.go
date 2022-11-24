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

const (
	errMsgDefault          = "Error"
	errMsgUnauthorized     = "Unauthorized"
	errMsgInvalidLoginData = "Invalid login data"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key            = []byte("super-secret-key")
	store          = sessions.NewCookieStore(key) // TODO
	authCookieName = "exampleproject_authenticated"
)

func applyAuthRequirement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, authCookieName)
		if err != nil {
			log.Log(log.GetFromSession, log.Err(err))
			http.Error(w, errMsgDefault, http.StatusInternalServerError)
			return
		}

		sessionValue := session.Values["authenticated"]
		if auth, ok := sessionValue.(bool); !ok {
			log.Log(log.CastInterface, log.UnexpectedType(sessionValue, true))
			http.Error(w, errMsgDefault, http.StatusInternalServerError)
			return
		} else if !auth {
			http.Error(w, errMsgUnauthorized, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func secret(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "secret message")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, authCookieName)
	if err != nil {
		log.Log(log.GetFromSession, log.Err(err))
		http.Error(w, errMsgDefault, http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Log(log.ReadBody)
		http.Error(w, errMsgDefault, http.StatusInternalServerError)
		return
	}

	var loginData api.LoginData
	if err := json.Unmarshal(body, &loginData); err != nil {
		log.Log(log.Unmarshal, log.JSONData(body))
		http.Error(w, errMsgDefault, http.StatusInternalServerError)
		return
	}

	authenticated, err := api.AuthenticateUser(r.Context(), loginData)
	if err != nil {
		http.Error(w, errMsgDefault, http.StatusInternalServerError)
		return
	}

	if !authenticated {
		http.Error(w, errMsgInvalidLoginData, http.StatusUnauthorized)
		return
	}

	session.Values["authenticated"] = true
	if err := store.Save(r, w, session); err != nil {
		log.Log(log.SaveToSession, log.Err(err))
		http.Error(w, errMsgDefault, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, authCookieName)
	if err != nil {
		log.Log(log.GetFromSession, log.Err(err))
		http.Error(w, errMsgDefault, http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = false
	if err := store.Save(r, w, session); err != nil {
		log.Log(log.SaveToSession, log.Err(err))
		http.Error(w, errMsgDefault, http.StatusInternalServerError)
		return
	}
}
