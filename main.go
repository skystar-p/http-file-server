package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore
var conf *Config

func main() {
	// read configuration
	config, err := ParseConfig("config.json")
	conf = config
	if err != nil {
		panic(err)
	}

	// initialize in-memory session store
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	store = sessions.NewCookieStore(randBytes)

	// router
	r := mux.NewRouter()

	// static file serving
	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// webauthn
	r.HandleFunc("/register-challenge", WebAuthnRegisterChallengeHandler).
		Methods("POST").
		Schemes("https")
	r.HandleFunc("/register", WebAuthnRegistrationHandler).
		Methods("POST").
		Schemes("https")

	r.HandleFunc("/authenticate", WebAuthnAuthenticationHandler)

	// receive signal
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go http.ListenAndServe(":10101", r)
	<-sig

	fmt.Printf("Signal received. Exit...\n")
	close(sig)
}
