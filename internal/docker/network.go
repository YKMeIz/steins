package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"log"
)

const (
	steinsNetwork = "steins-network"
)

func isNetworkExisted() bool {
	networkResource, _ := cli.NetworkInspect(context.Background(), steinsNetwork, types.NetworkInspectOptions{
		Scope:   "local",
		Verbose: false,
	})

	if networkResource.Name == steinsNetwork {
		return true
	}

	return false
}

func removeNetwork() {
	if err := cli.NetworkRemove(context.Background(), steinsNetwork); err != nil {
		log.Fatalln(err)
	}
}

func createNetwork() error {
	_, err := cli.NetworkCreate(context.Background(), steinsNetwork, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		Scope:          "local",
		EnableIPv6:     false,
		IPAM:           nil,
		Internal:       false,
		Attachable:     false,
		Ingress:        false,
		ConfigOnly:     false,
		ConfigFrom:     nil,
		Options:        nil,
		Labels:         nil,
	})
	return err
}
