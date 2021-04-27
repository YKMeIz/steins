package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"strconv"
	"sync"
)

func GetVirtualHosts() (sync.Map, error) {
	m := sync.Map{}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return sync.Map{}, err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		domain, ok := container.Labels["proxy.steins.server.name"]
		if !ok {
			continue
		}

		port, ok := container.Labels["proxy.steins.server.proxy_port"]
		if !ok {
			port = strconv.Itoa(80)
		}

		ipAddr := container.NetworkSettings.Networks[container.HostConfig.NetworkMode].IPAddress

		m.Store(domain, ipAddr + ":" + port)
	}

	return m, nil
}
