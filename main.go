package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	githubactions "github.com/sethvargo/go-githubactions"
)

func getDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func doesImageExist(ctx context.Context, client *client.Client, image, authStr string) bool {
	_, err := client.DistributionInspect(ctx, image, authStr)

	return err == nil
}

func main() {
	ctx := context.Background()
	action := githubactions.New()

	config, err := NewFromInputs(action)
	if err != nil {
		action.Fatalf("Failed to get inputs: %s", err)
	}
	action.Debugf(config.String())

	client, err := getDockerClient()
	if err != nil {
		action.Fatalf("Failed to get docker client: %s", err)
	}

	authStr, err := config.GetAuthString()
	if err != nil {
		action.Fatalf("Failed to get auth string: %s", err)
	}

	imageExists := doesImageExist(ctx, client, config.image, authStr)

	client.Close()

	if imageExists {
		fmt.Println("Image exists")
		os.Exit(0)
	}

	action.Fatalf("Image does not exist")
}
