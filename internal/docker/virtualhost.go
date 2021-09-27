package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"strconv"
	"strings"
	"sync"
)

type VirtualHosts struct {
	sync.Map
}

func GetVirtualHosts() (*VirtualHosts, error) {
	m := VirtualHosts{}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
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

		m.Store(domain, ipAddr+":"+port)
	}

	return &m, nil
}

func (vh *VirtualHosts) LoadWithReacquire(key string) (string, bool) {
	ip, ok := vh.Load(strings.Split(key, ":")[0])

	if !ok {
		v, err := GetVirtualHosts()
		if err != nil {
			return "", false
		}
		v.Range(func(key, value interface{}) bool {
			vh.Store(key, value)
			return true
		})

		ip, ok = vh.Load(strings.Split(key, ":")[0])
	}

	return ip.(string), ok
}
