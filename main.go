package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const serverAddress = ":8080"

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

var monitors = []Monitor{
	{ID: "sellers", Name: "Sellers", URL: "https://sellers.importacionesamexico.com.mx/"},
	{ID: "friopuro", Name: "Frío Puro", URL: "https://friopuro.com.mx/"},
	{ID: "imxtime", Name: "IMXTime", URL: "https://imxtime.importacionesamexico.com.mx/"},
	{ID: "imx", Name: "Importaciones a México", URL: "https://www.importacionesamexico.com.mx/"},
	{ID: "cadebot", Name: "Cadebot", URL: "https://www.importacionesamexico.com.mx/cadebot"},
}

type Monitor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type CheckResult struct {
	ID           string    `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	URL          string    `json:"url"`
	Status       string    `json:"status"`
	HTTPStatus   string    `json:"http_status,omitempty"`
	StatusCode   int       `json:"status_code,omitempty"`
	ResponseTime int64     `json:"response_time_ms"`
	CheckedAt    time.Time `json:"checked_at"`
	Error        string    `json:"error,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func checkWebsite(ctx context.Context, monitor Monitor) CheckResult {
	start := time.Now()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, monitor.URL, nil)
	if err != nil {
		return CheckResult{
			ID:        monitor.ID,
			Name:      monitor.Name,
			URL:       monitor.URL,
			Status:    "OFFLINE",
			CheckedAt: time.Now(),
			Error:     err.Error(),
		}
	}

	response, err := httpClient.Do(request)
	duration := time.Since(start)
	checkedAt := time.Now()

	if err != nil {
		return CheckResult{
			ID:           monitor.ID,
			Name:         monitor.Name,
			URL:          monitor.URL,
			Status:       "OFFLINE",
			ResponseTime: duration.Milliseconds(),
			CheckedAt:    checkedAt,
			Error:        err.Error(),
		}
	}

	defer response.Body.Close()

	status := "OFFLINE"
	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusBadRequest {
		status = "ONLINE"
	}

	return CheckResult{
		ID:           monitor.ID,
		Name:         monitor.Name,
		URL:          monitor.URL,
		Status:       status,
		HTTPStatus:   response.Status,
		StatusCode:   response.StatusCode,
		ResponseTime: duration.Milliseconds(),
		CheckedAt:    checkedAt,
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

func checkAllMonitors(ctx context.Context) []CheckResult {
	results := make([]CheckResult, len(monitors))
	var waitGroup sync.WaitGroup

	for index, monitor := range monitors {
		waitGroup.Add(1)
		go func(index int, monitor Monitor) {
			defer waitGroup.Done()
			results[index] = checkWebsite(ctx, monitor)
		}(index, monitor)
	}

	waitGroup.Wait()
	return results
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		fmt.Println("Error encoding response:", err)
	}
}

func allowCORS(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{
		"http://localhost:4321": true,
		"http://127.0.0.1:4321": true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigins[r.Header.Get("Origin")] {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", http.MethodGet+", OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, HealthResponse{
			Status:  "ok",
			Service: "gopulse-monitor",
		})
	})

	mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")

		if rawURL == "" {
			http.Error(w, "url parameter is required", http.StatusBadRequest)
			return
		}

		if !isValidURL(rawURL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		writeJSON(w, checkWebsite(r.Context(), Monitor{URL: rawURL}))
	})

	mux.HandleFunc("/api/monitors", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, checkAllMonitors(r.Context()))
	})

	fmt.Println("GoPulse API running on http://localhost:8080")

	if err := http.ListenAndServe(serverAddress, allowCORS(mux)); err != nil {
		fmt.Println("Server error:", err)
	}
}
