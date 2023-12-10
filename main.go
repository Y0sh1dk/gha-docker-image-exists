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
		panic(err)
	}

	fmt.Println(config)

	client, err := getDockerClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	authStr, err := config.GetAuthString()
	if err != nil {
		panic(err)
	}

	imageExists := doesImageExist(ctx, client, config.image, authStr)

	if imageExists {
		fmt.Println("Image exists")
		os.Exit(0)
	}

	fmt.Println("Image does not exist")
	os.Exit(1)
}
