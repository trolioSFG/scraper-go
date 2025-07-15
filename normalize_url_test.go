package main

import (
	//	"log"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expectedURL string
	}{
		{
			name:        "remove scheme",
			inputURL:    "https://blog.boot.dev/path",
			expectedURL: "blog.boot.dev/path",
		},
		{
			name:        "remove trailing /",
			inputURL:    "https://blog.boot.dev/path/",
			expectedURL: "blog.boot.dev/path",
		},
		{
			name:        "remove alt scheme",
			inputURL:    "http://blog.boot.dev/path",
			expectedURL: "blog.boot.dev/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAILED: Unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expectedURL {
				t.Errorf("Test %v - %s FAIL: Expected %v vs actual %v",
					i, tc.name, tc.expectedURL, actual)
			}
		})
	}
}
