package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/duo-labs/webauthn/protocol"
)

func WebAuthnRegisterChallengeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("starting register\n")
	user := MyUser{id: make([]byte, 16)}
	authSelect := protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.AuthenticatorAttachment(protocol.CrossPlatform),
		RequireResidentKey: false,
		UserVerification: protocol.VerificationPreferred,
	}
	conveyencePref := protocol.ConveyancePreference(protocol.PreferNoAttestation)
	options, sessionData, err := web.BeginRegistration(&user,
		webauthn.WithAuthenticatorSelection(authSelect),
		webauthn.WithConveyancePreference(conveyencePref))
	if err != nil {
		http.Error(w, "error when beginning webauthn registration", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	session, err := store.Get(r, "webauthn")
	if err != nil {
		http.Error(w, "unable to get session object", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	session.Values["registration-data"] = &sessionData
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(options)
	if err != nil {
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
}

func WebAuthnRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	user := MyUser{id: make([]byte, 16)}
	session, err := store.Get(r, "webauthn")
	if err != nil {
		http.Error(w, "unable to get session object", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	sessionData, ok := session.Values["registration-data"].(webauthn.SessionData)
	if !ok {
		http.Error(w, "unable to get session data", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	credential, err := web.FinishRegistration(&user, sessionData, r)
	if err != nil {
		http.Error(w, "unable to finish registration", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	session.Values["credential"] = credential
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	fmt.Fprintf(w, "Registration success")
}

func WebAuthnAuthenticateChallengeHandler(w http.ResponseWriter, r *http.Request) {
	user := MyUser{id: make([]byte, 16)}
	options, sessionData, err := web.BeginLogin(&user)
	if err != nil {
		http.Error(w, "unable to begin webauthn login", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	session, err := store.Get(r, "webauthn")
	if err != nil {
		http.Error(w, "unable to get session object", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	session.Values["authentication-data"] = &sessionData
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(options)
	if err != nil {
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
}

func WebAuthnAuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	user := MyUser{id: make([]byte, 16)}
	session, err := store.Get(r, "webauthn")
	if err != nil {
		http.Error(w, "unable to get session object", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	sessionData, ok := session.Values["authentication-data"].(webauthn.SessionData)
	if !ok {
		http.Error(w, "unable to get session data", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	credential, err := web.FinishLogin(&user, sessionData, r)
	if err != nil {
		http.Error(w, "unable to finish login", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	if credential.Authenticator.CloneWarning {
		http.Error(w, "credential clone warning", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "credential clone warning\n")
		return
	}

	session.Values["credential"] = credential
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	fmt.Fprintf(w, "Authentication success!")
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
