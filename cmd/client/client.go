package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var target = flag.String("target", "http://localhost:8090", "request target")

func main() {
	flag.Parse()
	client := new(http.Client)
	client.Timeout = 10 * time.Second

	for range time.Tick(1 * time.Second) {
		// resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", *target))
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/some-data", *target), nil)
		req.Header.Add("X-Forwarded-For", "217.115.97.69")
		resp, err := client.Do(req)
		if err == nil {
			log.Printf("response %d", resp.StatusCode)
			log.Printf(resp.Header.Get("lb-from"))
		} else {
			log.Printf("error %s", err)
		}
	}
}
