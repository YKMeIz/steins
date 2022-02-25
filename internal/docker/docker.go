package docker

import (
	"github.com/docker/docker/client"
	"log"
)

var (
	cli *client.Client
)

func NetworkInit() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalln(err)
	}

	if !isNetworkExisted() {
		if err := createNetwork(); err != nil {
			log.Fatalln(err)
		}
	}
}
