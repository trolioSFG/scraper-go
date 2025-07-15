package main

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	rsp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	if rsp.StatusCode >= 400 {
		return "", fmt.Errorf("Bad request: %d", rsp.StatusCode)
	}

	if !strings.HasPrefix(rsp.Header["Content-Type"][0], "text/html") {
		/**
		for k, v := range rsp.Header {
			fmt.Printf("%v: %v\n", k, v)
		}
		**/
		return "", fmt.Errorf("No html content")
	}

	buf, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v", err)
	}

	defer rsp.Body.Close()
	return string(buf), nil
}

func crawlPage(rawBaseURL, currentURL string, pages map[string]int) {
	if !strings.HasPrefix(currentURL, rawBaseURL) {
		return
	}

	normURL, err := normalizeURL(currentURL)
	if err != nil {
		return
	}

	_, ok := pages[normURL]
	if ok {
		pages[normURL] += 1
		return
	}

	pages[normURL] = 1
	fmt.Printf("Crawling %s\n", currentURL)

	html, err := getHTML(currentURL)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
	
	// fmt.Printf("%s\n", html)
	links, err := getURLSFromHTML(html, rawBaseURL)
	for _, l := range links {
		crawlPage(rawBaseURL, l, pages)
	}
}
		
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("no website provided\n")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Printf("too many arguments provided\n")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %v\n", os.Args[1])

	pages := make(map[string]int)
	crawlPage(os.Args[1], os.Args[1], pages)

	fmt.Printf("================================\n")
	for k, v := range pages {
		fmt.Printf("%s: %d\n", k, v)
	}
}

