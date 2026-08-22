package main

import (
	"fmt"
	"net/http"
)

func main() {
	url := "https://example.com"

	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer response.Body.Close()

	fmt.Println("URL:", url)
	fmt.Println("Status:", response.Status)
	fmt.Println("Status code:", response.StatusCode)
}
