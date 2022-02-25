package podman

import (
	"log"
	"os"
)

const (
	socketRoot = "/run/podman/podman.sock"
)

func getSock() string {
	if _, err := os.Stat("/run/podman/podman.sock"); err == nil {
		log.Println("detect podman.sock under root user")
		return "unix:" + socketRoot
	}

	socketUser := "unix:" + os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock"
	if _, err := os.Stat(socketUser); err == nil {
		log.Println("detect podman.sock under " + socketUser + " user")
		return socketUser
	}

	log.Println("podman.sock is not activated")
	return ""
}
