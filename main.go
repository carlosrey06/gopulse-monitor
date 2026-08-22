package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func checkWebsite(url string) {
	start := time.Now()

	response, err := http.Get(url)

	duration := time.Since(start)

	if err != nil {
		fmt.Println("Status: OFFLINE")
		fmt.Println("Error:", err)
		fmt.Println("Response time:", duration.Milliseconds(), "ms")
		return
	}

	defer response.Body.Close()

	status := "ONLINE"

	if response.StatusCode >= 400 {
		status = "OFFLINE"
	}

	fmt.Println("URL:", url)
	fmt.Println("Status:", status)
	fmt.Println("HTTP:", response.Status)
	fmt.Println("Status code:", response.StatusCode)
	fmt.Println("Response time:", duration.Milliseconds(), "ms")
}
func main() {
	websites := []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
	}

	start := time.Now()

	var wg sync.WaitGroup

	for _, website := range websites {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			checkWebsite(url)
		}(website)
	}

	wg.Wait()

	totalDuration := time.Since(start)

	fmt.Println("Total time:", totalDuration.Milliseconds(), "ms")
}
