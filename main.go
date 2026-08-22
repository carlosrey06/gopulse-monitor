package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CheckResult struct {
	URL          string
	Status       string
	HTTPStatus   string
	StatusCode   int
	ResponseTime int64
	Error        string
}

func checkWebsite(url string) CheckResult {
	start := time.Now()

	response, err := http.Get(url)

	duration := time.Since(start)

	if err != nil {
		return CheckResult{
			URL:          url,
			Status:       "OFFLINE",
			ResponseTime: duration.Milliseconds(),
			Error:        err.Error(),
		}
	}

	defer response.Body.Close()

	status := "ONLINE"

	if response.StatusCode >= 400 {
		status = "OFFLINE"
	}

	return CheckResult{
		URL:          url,
		Status:       status,
		HTTPStatus:   response.Status,
		StatusCode:   response.StatusCode,
		ResponseTime: duration.Milliseconds(),
	}
}

func printResult(result CheckResult) {
	fmt.Println("URL:", result.URL)
	fmt.Println("Status:", result.Status)

	if result.Error != "" {
		fmt.Println("Error:", result.Error)
	} else {
		fmt.Println("HTTP:", result.HTTPStatus)
		fmt.Println("Status code:", result.StatusCode)
	}

	fmt.Println("Response time:", result.ResponseTime, "ms")
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

			result := checkWebsite(url)
			printResult(result)
		}(website)
	}

	wg.Wait()

	totalDuration := time.Since(start)

	fmt.Println("Total time:", totalDuration.Milliseconds(), "ms")
}
