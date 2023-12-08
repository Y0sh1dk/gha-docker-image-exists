package main

import (
	"context"
	"os"

	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

func getEnvDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getAuthString(username, password, serverAddress string) (string, error) {
	authConfig := registry.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: serverAddress,
	}

	authStr, err := registry.EncodeAuthConfig(authConfig)
	if err != nil {
		return "", err
	}
	return authStr, nil
}

func getDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func doesImageExist(ctx context.Context, client *client.Client, image, authStr string) bool {
	_, err := client.DistributionInspect(ctx, image, authStr)
	return err == nil
}

func main() {
	ctx := context.Background()

	username := getEnvDefault("DOCKER_USERNAME", "")
	password := getEnvDefault("DOCKER_PASSWORD", "")
	serverAddress := getEnvDefault("DOCKER_SERVER_ADDRESS", "https://index.docker.io/v1/")
	image, exists := os.LookupEnv("DOCKER_IMAGE")
	if !exists {
		panic("DOCKER_IMAGE environment variable not set")
	}

	client, err := getDockerClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	authStr, err := getAuthString(username, password, serverAddress)
	if err != nil {
		panic(err)
	}

	imageExists := doesImageExist(ctx, client, image, authStr)

	if imageExists {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
