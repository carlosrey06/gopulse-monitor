package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func checkWebsite(rawURL string) CheckResult {
	start := time.Now()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Get(rawURL)

	duration := time.Since(start)

	if err != nil {
		return CheckResult{
			URL:          rawURL,
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
		URL:          rawURL,
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

func isValidURL(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)

	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return parsedURL.Host != ""
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := HealthResponse{
			Status:  "ok",
			Service: "gopulse-monitor",
		}

		err := json.NewEncoder(w).Encode(response)

		if err != nil {
			fmt.Println("Error encoding response:", err)
		}
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")

		if rawURL == "" {
			http.Error(w, "url parameter is required", http.StatusBadRequest)
			return
		}

		if !isValidURL(rawURL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		result := checkWebsite(rawURL)

		err := json.NewEncoder(w).Encode(result)

		if err != nil {
			fmt.Println("Error encoding response:", err)
		}
	})

	fmt.Println("GoPulse API running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
