package docker

import (
	"testing"
)

func TestNetworkInit(t *testing.T) {
	if isNetworkExisted() {
		removeNetwork()
	}
}

func TestNetworkCreate(t *testing.T) {
	if err := createNetwork(); err != nil {
		t.Error(err)
	}
}

func TestNetworkCheck(t *testing.T) {
	if !isNetworkExisted() {
		t.Error("local network created but not found\n")
	}
}

func TestNetworkCleanup(t *testing.T) {
	removeNetwork()
}
