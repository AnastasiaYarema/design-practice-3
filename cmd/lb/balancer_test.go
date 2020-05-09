package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetServerURLByClientAddressHashing1(t *testing.T) {
	url, _ := getServerURLByClientAddressHashing(map[string]bool {
		"server1:8080": true,
		"server2:8080": true,
		"server3:8080": true,
	}, "217.115.97.69")

	assert.Equal(t, url, "server1:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashing2(t *testing.T) {
	url, _ := getServerURLByClientAddressHashing(map[string]bool {
		"server1:8080": true,
		"server2:8080": true,
		"server3:8080": true,
	}, "217.115.97.68")

	assert.Equal(t, url, "server2:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashing3(t *testing.T) {
	url, _ := getServerURLByClientAddressHashing(map[string]bool {
		"server1:8080": true,
		"server2:8080": true,
		"server3:8080": true,
	}, "217.115.97.71")

	assert.Equal(t, url, "server3:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashingWithOneDisabledServer(t *testing.T) {
	url, _ := getServerURLByClientAddressHashing(map[string]bool {
		"server1:8080": false,
		"server2:8080": true,
		"server3:8080": true,
	}, "217.115.97.68")

	assert.Equal(t, url, "server2:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashingWithTwoDisabledServer1(t *testing.T) {
	url, _ := getServerURLByClientAddressHashing(map[string]bool {
		"server1:8080": false,
		"server2:8080": false,
		"server3:8080": true,
	}, "217.115.97.68")

	assert.Equal(t, url, "server3:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashingWithTwoDisabledServer2(t *testing.T) {
	url, _ := getServerURLByClientAddressHashing(map[string]bool {
		"server1:8080": true,
		"server2:8080": false,
		"server3:8080": false,
	}, "217.115.97.68")

	assert.Equal(t, url, "server1:8080", "Wrong url")
}

func TestGetServerURLByClientAddressHashingAllDisabledServers(t *testing.T) {
	_, err := getServerURLByClientAddressHashing(map[string]bool {
		"server1:8080": false,
		"server2:8080": false,
		"server3:8080": false,
	}, "217.115.97.68")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "There are no available servers")
}

func TestGetServerURLByClientAddressHashingWithEmptyClientAddressString(t *testing.T) {
	_, err := getServerURLByClientAddressHashing(map[string]bool {
		"server1:8080": false,
		"server2:8080": false,
		"server3:8080": false,
	}, "")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Client address cannot be empty")
}