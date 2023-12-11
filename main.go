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

	action.Debugf("Getting inputs...")
	config, err := NewFromInputs(action)
	if err != nil {
		panic(err)
	}
	action.Debugf("Got inputs")
	fmt.Println(config)

	action.Debugf("Getting docker client...")
	client, err := getDockerClient()
	if err != nil {
		panic(err)
	}
	action.Debugf("Got docker client")

	action.Debugf("Getting auth string...")
	authStr, err := config.GetAuthString()
	if err != nil {
		panic(err)
	}
	action.Debugf("Got auth string")

	action.Debugf("Checking if image exists...")
	imageExists := doesImageExist(ctx, client, config.image, authStr)
	action.Debugf("Checked if image exists")

	client.Close()

	if imageExists {
		fmt.Println("Image exists")
		os.Exit(0)
	}

	fmt.Println("Image does not exist")
	os.Exit(1)
}
