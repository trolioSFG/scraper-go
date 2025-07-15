package main

import (
	"strings"
	"net/url"
)

func normalizeURL(raw string) (string, error) {
	// parts := strings.Split(url, "://")
	// parts[1] = strings.TrimSuffix(parts[1], "/")
	parsedURL, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	return parsedURL.Host + strings.TrimSuffix(parsedURL.Path, "/"), nil
}
