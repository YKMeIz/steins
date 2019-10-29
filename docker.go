package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"strings"
	"sync"
)

func getVirtualHosts() (sync.Map, error) {
	m := sync.Map{}
	cli, err := client.NewEnvClient()
	if err != nil {
		return sync.Map{}, err
	}

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return sync.Map{}, err
	}

	// Combine default bridge network and self created networks.
	// Only support bridge driver.
	var bridgeNetworks []types.NetworkResource
	for _, v := range networks {
		if v.Driver == "bridge" {
			bridgeNetworks = append(bridgeNetworks, v)
		}
	}

	for _, network := range bridgeNetworks {
		for _, container := range network.Containers {
			info, err := cli.ContainerInspect(context.Background(), container.Name)
			if err != nil {
				log.Println(err)
				continue
			}
			domain, ok := info.Config.Labels["proxy.steins.server.name"]
			if !ok {
				continue
			}
			port, ok := info.Config.Labels["proxy.steins.server.proxy_port"]
			if !ok {
				port = strconv.Itoa(80)
			}

			for _, v := range strings.Split(domain, ",") {
				host := ":" + port
				if network.Name == "bridge" {
					host = info.NetworkSettings.Networks["bridge"].IPAddress + host
				} else {
					host = info.Name[1:] + "." + network.Name + host
				}
				m.Store(v, host)
			}
		}
	}

	return m, nil
}
