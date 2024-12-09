package api

import "github.com/go-webauthn/webauthn/webauthn"

type State struct {
	webAuthn     *webauthn.WebAuthn
	users        map[string]*User
	sessionStore map[string]*webauthn.SessionData
}

func NewState(webAuthn *webauthn.WebAuthn) *State {
	return &State{
		webAuthn:     webAuthn,
		sessionStore: make(map[string]*webauthn.SessionData),
		users:        make(map[string]*User),
	}
}
