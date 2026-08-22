package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	url := "https://example.com"

	start := time.Now()

	response, err := http.Get(url)

	duration := time.Since(start)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer response.Body.Close()

	fmt.Println("URL:", url)
	fmt.Println("Status:", response.Status)
	fmt.Println("Status code:", response.StatusCode)
	fmt.Println("Response time:", duration.Milliseconds(), "ms")
}
