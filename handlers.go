package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/duo-labs/webauthn/webauthn"
)

func WebAuthnRegisterChallengeHandler(w http.ResponseWriter, r *http.Request) {
	user := MyUser{}
	options, sessionData, err := web.BeginRegistration(&user)
	if err != nil {
		http.Error(w, "error when beginning webauthn registration", http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "webauthn")
	if err != nil {
		http.Error(w, "unable to get session object", http.StatusInternalServerError)
		return
	}

	session.Values["registration-data"] = sessionData
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}

	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(options)
	if err != nil {
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		return
	}
}

func WebAuthnRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	user := MyUser{}
	session, err := store.Get(r, "webauthn")
	if err != nil {
		http.Error(w, "unable to get session object", http.StatusInternalServerError)
		return
	}

	sessionData, ok := session.Values["registration-data"].(webauthn.SessionData)
	if !ok {
		http.Error(w, "unable to get session data", http.StatusInternalServerError)
		return
	}

	credential, err := web.FinishRegistration(&user, sessionData, r)
	fmt.Printf("credential: %+v\n", credential)
	if err != nil {
		http.Error(w, "unable to finish registration", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Registration success")
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
