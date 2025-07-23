package common

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func PingServerLoop(url string) {
	if url == "" {
		fmt.Println("PingServerLoop: URL is empty, skipping ping.")
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		delay := r.Intn(5) + 10
		fmt.Printf("Waiting %d minutes before pinging %s...\n", delay, url)
		time.Sleep(time.Duration(delay) * time.Minute)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("PingServerLoop: Error pinging %s: %v\n", url, err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("PingServerLoop: Successfully pinged %s, status code: %d\n", url, resp.StatusCode)
		} else {
			fmt.Printf("PingServerLoop: Received non-OK status code %d from %s\n", resp.StatusCode, url)
		}
	}
}