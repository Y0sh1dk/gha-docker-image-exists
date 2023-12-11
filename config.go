package main

import (
	"fmt"

	"github.com/docker/docker/api/types/registry"
	"github.com/sethvargo/go-githubactions"
)

type Config struct {
	username      string
	password      string
	serverAddress string
	image         string
	action        *githubactions.Action
}

func (c Config) String() string {
	return fmt.Sprintf("Config{username: %s, password: #####, serverAddress: %s, image: %s}", c.username, c.serverAddress, c.image)
}

func (c *Config) GetAuthString() (string, error) {
	c.action.Debugf("Getting auth string...")
	defer c.action.Debugf("Got auth string")

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
	action.Debugf("Getting inputs...")
	defer action.Debugf("Got inputs")

	username := getInputDefault(action, "username", "")
	password := getInputDefault(action, "password", "")
	serverAddress := getInputDefault(action, "server_address", "")
	image := getInputDefault(action, "image", "")

	c := Config{
		username:      username,
		password:      password,
		serverAddress: serverAddress,
		image:         image,
		action:        action,
	}

	return &c, nil
}

func getInputDefault(action *githubactions.Action, key, fallback string) string {
	action.Debugf("Getting input %s...", key)
	defer action.Debugf("Got input %s", key)

	if value := action.GetInput(key); value != "" {
		return value
	}

	return fallback
}
