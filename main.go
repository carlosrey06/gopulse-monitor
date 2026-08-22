package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CheckResult struct {
	URL          string `json:"url"`
	Status       string `json:"status"`
	HTTPStatus   string `json:"http_status,omitempty"`
	StatusCode   int    `json:"status_code,omitempty"`
	ResponseTime int64  `json:"response_time_ms"`
	Error        string `json:"error,omitempty"`
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

	results := make(chan CheckResult)

	for _, website := range websites {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			result := checkWebsite(url)
			results <- result
		}(website)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var checkResults []CheckResult

	for result := range results {
		checkResults = append(checkResults, result)
	}

	for _, result := range checkResults {
		printResult(result)
	}

	jsonData, err := json.MarshalIndent(checkResults, "", "  ")

	if err != nil {
		fmt.Println("Error generating JSON:", err)
		return
	}

	fmt.Println()
	fmt.Println("JSON:")
	fmt.Println(string(jsonData))

	totalDuration := time.Since(start)

	fmt.Println("Total time:", totalDuration.Milliseconds(), "ms")
}
