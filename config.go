package main

import (
	"encoding/json"
	"errors"
	"os"
)

type WebAuthnConfig struct {
	RPID string `json:"rpId"`
	RPDisplayName string `json:"rpDisplayName"`
	RPOrigin string `json:"rpOrigin"`
	RPIcon string `json:"rpIcon"`
}

type Config struct {
	RootPath string `json:"rootPath"`
	WebAuthn WebAuthnConfig `json:"webauthn"`
}

func ParseConfig(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	defer f.Close()
	if err != nil {
		return nil, errors.New("error when opening config file")
	}
	jsonDecoder := json.NewDecoder(f)
	config := Config{}

	err = jsonDecoder.Decode(&config)
	if err != nil {
		return nil, errors.New("error when reading config file")
	}

	return &config, nil
}
