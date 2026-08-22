package main

import (
	"fmt"
	"net/http"
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
	checkWebsite("https://example.com")
}
