package main

import (
	"github.com/docker/docker/api/types/registry"
	"github.com/sethvargo/go-githubactions"
)

type Config struct {
	username      string
	password      string
	serverAddress string
	image         string
}

func (c *Config) GetAuthString() (string, error) {
	authConfig := registry.AuthConfig{
		Username:      c.username,
		Password:      c.password,
		ServerAddress: c.serverAddress,
	}

	authStr, err := registry.EncodeAuthConfig(authConfig)
	if err != nil {
		return "", err
	}
	return authStr, nil
}

func NewFromInputs(action *githubactions.Action) (*Config, error) {
	username := getInputDefault(action, "username", "")
	password := getInputDefault(action, "password", "")
	serverAddress := getInputDefault(action, "server_address", "")
	image := getInputDefault(action, "image", "")

	c := Config{
		username:      username,
		password:      password,
		serverAddress: serverAddress,
		image:         image,
	}
	return &c, nil
}

func getInputDefault(action *githubactions.Action, key, fallback string) string {
	if value := action.GetInput(key); value == "" {
		return value
	}
	return fallback
}
