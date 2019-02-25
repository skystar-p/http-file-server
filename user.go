package main

import (
	"github.com/duo-labs/webauthn/webauthn"
)

type MyUser struct {
	id []byte
}

func (user *MyUser) WebAuthnID() []byte {
	return user.id
}

func (user *MyUser) WebAuthnName() string {
	return "skystar"
}

func (user *MyUser) WebAuthnDisplayName() string {
	return "skystar-dn"
}

func (user *MyUser) WebAuthnIcon() string {
	return ""
}

func (user *MyUser) WebAuthnCredentials() []webauthn.Credential {
	return []webauthn.Credential{}
}
