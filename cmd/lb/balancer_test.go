package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetServerURLByClientAddressHashing1(t *testing.T) {
    availableServers = []string {
		"server1:8080",
		"server2:8080",
		"server3:8080",
    }
	url, _ := getServerURLByClientAddressHashing("217.115.97.69")

	assert.Equal(t, url, "server1:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashing2(t *testing.T) {
    availableServers = []string {
		"server1:8080",
		"server2:8080",
		"server3:8080",
    }
	url, _ := getServerURLByClientAddressHashing("217.115.97.68")

	assert.Equal(t, url, "server2:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashing3(t *testing.T) {
    availableServers = []string {
		"server1:8080",
		"server2:8080",
		"server3:8080",
    }
	url, _ := getServerURLByClientAddressHashing("217.115.97.71")

	assert.Equal(t, url, "server3:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashingWithOneDisabledServer(t *testing.T) {
    availableServers = []string {
		"server2:8080",
		"server3:8080",
    }
	url, _ := getServerURLByClientAddressHashing("217.115.97.68")

	assert.Equal(t, url, "server2:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashingWithTwoDisabledServer1(t *testing.T) {
    availableServers = []string {
		"server3:8080",
    }
	url, _ := getServerURLByClientAddressHashing("217.115.97.68")

	assert.Equal(t, url, "server3:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashingWithTwoDisabledServer2(t *testing.T) {
    availableServers = []string {
		"server1:8080",
    }
	url, _ := getServerURLByClientAddressHashing("217.115.97.68")

	assert.Equal(t, url, "server1:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashingAllDisabledServers(t *testing.T) {
    availableServers = []string {}
	_, err := getServerURLByClientAddressHashing("217.115.97.68")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "There are no available servers")
}

func TestGetServerURLByClientAddressHashingWithEmptyClientAddressString(t *testing.T) {
    availableServers = []string {}
	_, err := getServerURLByClientAddressHashing("")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Client address cannot be empty")
}
