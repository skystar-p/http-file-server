package main

import (
	"fmt"
	"os"
	"encoding/json"
	"encoding/base64"
	"crypto/sha256"
)

type WebAuthnRegistrationChallenge struct {
	RPID string `json:"rpid"`
	UserName string `json:"username"`
	DisplayName string `json:"displayName"`
	Timeout int `json:"timeout"`
	Challenge string `json:"challenge"`
}

type WebAuthnAuthenticationChallenge struct {
	AllowCredentials CredentialDescriptor `json:"allowCredentials"`
	Timeout int `json:"timeout"`
	Challenge string `json:"challenge"`
}

type CredentialDescriptor struct {
	ID string `json:"id"`
	Transports []string `json:"transports"`
	Type string `json:"type"`
}

type WebAuthnRegistrationCredential struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Response AuthenticatorAttestationResponse `json:"response"`
}

type WebAuthnAuthenticationCredential struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Response AuthenticatorAssertionResponse `json:"response"`
}

type AuthenticatorAttestationResponse struct {
	AttestationObj string `json:"attestationObject"`
	ClientDataJSON string `json:"clientDataJSON"`
}

type AuthenticatorAssertionResponse struct {
	AuthenticationData string `json:"authenticationData"`
	ClientDataJSON string `json:"clientDataJSON"`
	Signature string `json:"signature"`
	UserHandle ClientData `json:"userHandle"`
}

type ClientData struct {
	Challenge string `json:"challenge"`
	HashAlgorithm string `json:"hashAlgorithm"`
	Origin string `json:"origin"`
	Type string `json:"type"`
}

const TYPE_CREATE string = "webauthn.create"
const TYPE_GET string = "webauthn.get"

func VerifyRegistrationRequest(req *WebAuthnRegistrationCredential, chal string) bool {
	fmt.Printf("starting registration request verification for %s\n", req.ID)

	// check type
	if req.Type != TYPE_CREATE {
		fmt.Fprintf(os.Stderr, "error: type is not TYPE_CREATE\n")
		return false
	}

	// unmarshal clientDataJSON
	clientData := ClientData{}
	clientDataByte := []byte(req.Response.ClientDataJSON)
	err := json.Unmarshal(clientDataByte, &clientData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to parse clientDataJSON\n")
		return false
	}

	if clientData.HashAlgorithm != chal {
		fmt.Fprintf(os.Stderr, "error: challenge mismatch\n")
		return false
	}

	// tokenBinding.status matching process skipped

	// compute sha-256 hash of clientDataJSON
	clientDataHash := sha256.Sum256(clientDataByte)

	// cbor decoder
	attestationObj, err := base64.StdEncoding.DecodeString(req.Response.AttestationObj)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: fail to b64decode attestationObj\n")
		return false
	}
	// err := cbor.Loads(attestationObj)
}
