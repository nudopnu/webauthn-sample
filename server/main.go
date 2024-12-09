package main

import (
	"fmt"
	"net/http"

	"githbu.com/nudopnu/webauthn-sample/internal/api"
	"github.com/go-webauthn/webauthn/webauthn"
	"vitess.io/vitess/go/vt/log"
)

func main() {
	wconfig := &webauthn.Config{
		RPDisplayName: "Go Webauthn",                // Display Name for your site
		RPID:          "localhost",                  // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost"}, // The origin URLs allowed for WebAuthn requests
	}

	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		log.Fatal(err)
	}
	state := api.NewState(webAuthn)

	// Handlers
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register/start", state.BeginRegistration)
	mux.HandleFunc("POST /register/finish", state.FinishRegistration)
	mux.HandleFunc("POST /login/start", state.BeginLogin)
	mux.HandleFunc("POST /login/finish", state.FinishLogin)

	fmt.Printf("Running server on http://0.0.0.0:8080\n")
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: api.HandleCORS(mux),
	}
	server.ListenAndServe()
}
