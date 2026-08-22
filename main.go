package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
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
		w.Header().Set("Content-Type", "application/json")

		url := r.URL.Query().Get("url")

		if url == "" {
			http.Error(w, "url parameter is required", http.StatusBadRequest)
			return
		}

		result := checkWebsite(url)

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
