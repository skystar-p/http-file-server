package main

import (
	// "fmt"
	"net/http"
	"encoding/json"
	"crypto/rand"
	"encoding/base64"
	// "github.com/gorilla/mux"
)

type WebAuthnRegistrationResponse struct {
	RPID string
	UserName string
	DisplayName string
	Timeout int
	Challenge string
}

func WebAuthnRegisterChallengeHandler(w http.ResponseWriter, r *http.Request) {
	chal := make([]byte, 64)
	_, err := rand.Read(chal)
	if err != nil {
		http.Error(w, "unable to generate random challenge", http.StatusInternalServerError)
		return
	}

	b64enc := base64.NewEncoding("utf-8")
	chalStr := b64enc.EncodeToString(chal)
	resp := WebAuthnRegistrationResponse{
		RPID: conf.WebAuthn.RPID,
		UserName: conf.WebAuthn.UserName,
		DisplayName: conf.WebAuthn.DisplayName,
		Timeout: conf.WebAuthn.Timeout,
		Challenge: chalStr,
	}

	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(resp)
	if err != nil {
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
	}
}

func WebAuthnRegistrationHandler(w http.ResponseWriter, r *http.Request) {

}

func WebAuthnAuthenticationHandler(w http.ResponseWriter, r *http.Request) {

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if false {
			// perform authentication on here...
			http.Error(w, "authentication error", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
