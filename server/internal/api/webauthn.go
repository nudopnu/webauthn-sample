package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
)

func (state *State) BeginRegistration(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// Create a new user or retrieve an existing one
	user, ok := state.users[username]
	if !ok {
		user = &User{
			ID:          []byte(username),
			Name:        username,
			DisplayName: username,
			Credentials: make([]webauthn.Credential, 0),
		}
		state.users[username] = user
	}

	// Start the registration
	options, session, err := state.webAuthn.BeginRegistration(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Save session data (in practice, use a secure store)
	state.sessionStore[username] = session

	// Send the registration options to the client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

func (state *State) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// Get the user
	user, ok := state.users[username]
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Retrieve session data
	session, ok := state.sessionStore[username]
	if !ok {
		http.Error(w, "session data not found", http.StatusBadRequest)
		return
	}

	// Parse the credential creation response
	credential, err := state.webAuthn.FinishRegistration(user, *session, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	// If creation was successful, store the credential object
	user.AddCredential(*credential)

	// TODO: add access + refresh token
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registration Success"))
}

func (state *State) BeginLogin(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	user, ok := state.users[username]
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Start the login
	options, sessionData, err := state.webAuthn.BeginLogin(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to start login: %v", err), http.StatusInternalServerError)
		return
	}

	// Save session data
	state.sessionStore[username] = sessionData

	// Send the login options to the client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

func (state *State) FinishLogin(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	user, ok := state.users[username]
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	// Retrieve session data
	sessionData, ok := state.sessionStore[username]
	if !ok {
		http.Error(w, "session data not found", http.StatusBadRequest)
		return
	}

	// Verify the assertion
	_, err := state.webAuthn.FinishLogin(user, *sessionData, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to finish login: %v", err), http.StatusInternalServerError)
		return
	}

	// TODO: add access + refresh token
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login Success"))
}
