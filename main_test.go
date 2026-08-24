package main

import "testing"

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		valid bool
	}{
		{name: "valid HTTPS URL", url: "https://example.com", valid: true},
		{name: "valid HTTP URL", url: "http://example.com", valid: true},
		{name: "plain text", url: "hola", valid: false},
		{name: "FTP URL", url: "ftp://example.com", valid: false},
		{name: "JavaScript URL", url: "javascript:alert(1)", valid: false},
		{name: "empty URL", url: "", valid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isValidURL(test.url); got != test.valid {
				t.Fatalf("isValidURL(%q) = %t, want %t", test.url, got, test.valid)
			}
		})
	}
}
