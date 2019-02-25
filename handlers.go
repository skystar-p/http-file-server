package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"crypto/rand"
	"encoding/base64"
)

func WebAuthnRegisterChallengeHandler(w http.ResponseWriter, r *http.Request) {
	chal := make([]byte, 64)
	_, err := rand.Read(chal)
	if err != nil {
		http.Error(w, "unable to generate random challenge", http.StatusInternalServerError)
		return
	}

	b64enc := base64.NewEncoding("utf-8")
	chalStr := b64enc.EncodeToString(chal)
	resp := WebAuthnRegistrationChallenge{
		RPID: conf.WebAuthn.RPID,
		UserName: conf.WebAuthn.UserName,
		DisplayName: conf.WebAuthn.DisplayName,
		Timeout: conf.WebAuthn.Timeout,
		Challenge: chalStr,
	}

	session, err := store.Get(r, "webauthn")
	if err != nil {
		http.Error(w, "unable to get session object", http.StatusInternalServerError)
		return
	}

	session.Values["chalStr"] = chalStr
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}

	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(resp)
	if err != nil {
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		return
	}
}

func WebAuthnRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	jsonDecoder := json.NewDecoder(r.Body)
	regCred := WebAuthnRegistrationCredential{}
	err := jsonDecoder.Decode(&regCred)
	if err != nil {
		http.Error(w, "unable to decode registration request", http.StatusBadRequest)
		return
	}

	// verify at here
}

func WebAuthnAuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	jsonDecoder := json.NewDecoder(r.Body)
	authCred := WebAuthnAuthenticationCredential{}
	err := jsonDecoder.Decode(&authCred)
	if err != nil {
		http.Error(w, "unable to decode authentication request", http.StatusBadRequest)
		return
	}
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
