package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"hash/adler32"
	"errors"
	"sort"

	"github.com/AnastasiaYarema/design-practice-3/httptools"
	"github.com/AnastasiaYarema/design-practice-3/signal"
)

var (
	port = flag.Int("port", 8090, "load balancer port")
	timeoutSec = flag.Int("timeout-sec", 3, "request timeout time in seconds")
	https = flag.Bool("https", false, "whether backends support HTTPs")

	traceEnabled = flag.Bool("trace", false, "whether to include tracing information into responses")
)

var (
	timeout = time.Duration(*timeoutSec) * time.Second
	serversPool = []string {
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}
	serversHealth = map[string]bool {
		"server1:8080": true,
		"server2:8080": true,
		"server3:8080": false,
	}
)

func scheme() string {
	if *https {
		return "https"
	}
	return "http"
}

func health(dst string) bool {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s://%s/health", scheme(), dst), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func forward(dst string, rw http.ResponseWriter, r *http.Request) error {
	ctx, _ := context.WithTimeout(r.Context(), timeout)
	fwdRequest := r.Clone(ctx)
	fwdRequest.RequestURI = ""
	fwdRequest.URL.Host = dst
	fwdRequest.URL.Scheme = scheme()
	fwdRequest.Host = dst

	resp, err := http.DefaultClient.Do(fwdRequest)
	if err == nil {
		for k, values := range resp.Header {
			for _, value := range values {
				rw.Header().Add(k, value)
			}
		}
		if *traceEnabled {
			rw.Header().Set("lb-from", dst)
		}
		log.Println("fwd", resp.StatusCode, resp.Request.URL)
		rw.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		_, err := io.Copy(rw, resp.Body)
		if err != nil {
			log.Printf("Failed to write response: %s", err)
		}
		return nil
	} else {
		log.Printf("Failed to get response from %s: %s", dst, err)
		rw.WriteHeader(http.StatusServiceUnavailable)
		return err
	}
}

// Calculates server url with algorithm of client address hashing considering servers' health
func getServerURLByClientAddressHashing(servers map[string]bool, clientAddress string) (string, error) {
	if (clientAddress == "") {
		return "", errors.New("Client address cannot be empty")
	}

	var availableServers []string

	for server, isAvailable := range servers {
		if isAvailable {
			availableServers = append(availableServers, server)
		}
	}

	sort.Strings(availableServers)

	if len(availableServers) == 0 {
		return "", errors.New("There are no available servers")
	}

	checksum := int(adler32.Checksum([]byte(clientAddress)))
	serverIndex := checksum % len(availableServers)
	targerServer := availableServers[serverIndex]

	return targerServer, nil
}

func main() {
	flag.Parse()

	for server := range serversHealth {
		server := server
		go func() {
			for range time.Tick(10 * time.Second) {
				log.Println(server, health(server))
				serversHealth[server] = health(server)
			}
		}()
	}

	frontend := httptools.CreateServer(*port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		targetServer, err := getServerURLByClientAddressHashing(serversHealth, r.RemoteAddr) // calculates target server url with algorithm of client address hashing
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusServiceUnavailable)
		} else {
			forward(targetServer, rw, r)
		}
	}))

	log.Println("Starting load balancer...")
	log.Printf("Tracing support enabled: %t", *traceEnabled)
	frontend.Start()
	signal.WaitForTerminationSignal()
}
