package main

import (
	"encoding/json"
	"errors"
	"os"
)

type WebAuthnConfig struct {
	RPID string `json:"rpId"`
	UserID string `json:"userId"`
	UserName string `json:"username"`
	DisplayName string `json:"displayName"`
	Timeout int `json:"timeout"`
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
