package docker

import (
	"github.com/docker/docker/client"
)

const (
	DOCKER_HOST = "unix:///home/razvan/.docker/desktop/docker.sock"
)

func CreateClient() (*client.Client, error) {

	client, err := client.NewClientWithOpts(client.WithHost(DOCKER_HOST))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client, nil
}
