package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"encoding/gob"

	"github.com/gorilla/mux"
	// "github.com/gorilla/sessions"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/boj/redistore"
)

var store *redistore.RediStore
var conf *Config
var web *webauthn.WebAuthn

func main() {
	// read configuration
	config, err := ParseConfig("config.json")
	conf = config
	if err != nil {
		panic(err)
	}

	// initialize webauthn
	web, err = webauthn.New(&webauthn.Config{
		RPID: conf.WebAuthn.RPID,
		RPDisplayName: conf.WebAuthn.RPDisplayName,
		RPOrigin: conf.WebAuthn.RPOrigin,
		RPIcon: conf.WebAuthn.RPIcon,
	})
	if err != nil {
		panic(err)
	}

	// initialize in-memory session store
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	// store = sessions.NewCookieStore(randBytes)
	store, err = redistore.NewRediStore(10, "tcp", ":6379", "", randBytes)
	if err != nil {
		panic(err)
	}
	gob.Register(webauthn.SessionData{})
	gob.Register(webauthn.Credential{})

	// router
	r := mux.NewRouter()

	// static file serving
	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// webauthn
	r.HandleFunc("/register-challenge", WebAuthnRegisterChallengeHandler).
		Methods("POST")
	r.HandleFunc("/register", WebAuthnRegistrationHandler).
		Methods("POST")

	r.HandleFunc("/authenticate-challenge", WebAuthnAuthenticateChallengeHandler).
		Methods("POST")
	r.HandleFunc("/authenticate", WebAuthnAuthenticationHandler).
		Methods("POST")

	// receive signal
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go http.ListenAndServe(":10101", r)
	<-sig

	fmt.Printf("Signal received. Exit...\n")
	close(sig)
}
