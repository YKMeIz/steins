package podman

import (
	"context"
	"errors"
	"github.com/YKMeIz/steins/internal/docker"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"strconv"
	"strings"
)

func GetVirtualHosts() (*docker.VirtualHosts, error) {
	sock := getSock()
	if sock == "" {
		return nil, errors.New("podman is not using")
	}

	connText, err := bindings.NewConnection(context.Background(), sock)
	if err != nil {
		return nil, err
	}

	opts := containers.ListOptions{
		Filters: make(map[string][]string),
	}
	opts.Filters["label"] = []string{"proxy.steins.server.name"}
	instances, err := containers.List(connText, &opts)
	if err != nil {
		return nil, err
	}

	m := docker.VirtualHosts{}

	for _, instance := range instances {
		domain, ok := instance.Labels["proxy.steins.server.name"]
		if !ok {
			continue
		}

		port, ok := instance.Labels["proxy.steins.server.proxy_port"]
		if !ok {
			port = strconv.Itoa(80)
		}

		instanceInfo, err := containers.Inspect(connText, instance.ID, nil)
		if err != nil {
			continue
		}

		ipAddr := instanceInfo.NetworkSettings.IPAddress
		m.Store(domain, ipAddr+":"+port)
	}

	return &m, nil
}

func LoadWithReacquire(vh *docker.VirtualHosts, key string) (string, bool) {
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
