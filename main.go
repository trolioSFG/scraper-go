package main

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"strings"
	"strconv"
	"sync"
)

type config struct {
	pages map[string]int
	baseURL string
	mu *sync.Mutex
	concurrencyControl chan struct{}
	wg *sync.WaitGroup
	maxPages int
	maxConcurrency int
}

const (
	CONCURRENCY_LEVEL int = 1	
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

func (cfg *config) crawlPage(currentURL string) {


	cfg.concurrencyControl <- struct{}{}
	defer func(){
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	cfg.mu.Lock()
	pagesSize := len(cfg.pages)
	cfg.mu.Unlock()

	if pagesSize >= cfg.maxPages {
		return
	}

	if !strings.HasPrefix(currentURL, cfg.baseURL) {
		fmt.Printf("Discarded URL: %s\n", currentURL)
		return
	}

	normURL, err := normalizeURL(currentURL)
	if err != nil {
		return
	}

	cfg.mu.Lock()
	_, ok := cfg.pages[normURL]
	if ok {
		cfg.pages[normURL] += 1
		// HERE: Forgot to unlock
		cfg.mu.Unlock()
		return
	}

	cfg.pages[normURL] = 1
	cfg.mu.Unlock()

	fmt.Printf("Crawling %s\n", currentURL)

	html, err := getHTML(currentURL)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}


	// fmt.Printf("%s\n", html)
	links, err := getURLSFromHTML(html, cfg.baseURL)
	if err != nil {
		return
	}

	for _, l := range links {
		cfg.wg.Add(1)
		go cfg.crawlPage(l)
	}


}
		
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("no website provided\n")
		os.Exit(1)
	}

	if len(os.Args) == 3 || len(os.Args) > 4 {
		fmt.Printf("Usage: %s <baseURL> <maxConcurrency> <maxPages>\n")
		os.Exit(1)
	}

	maxConcurrency := 5
	maxPages := 10
	var err error

	if len(os.Args) == 4 {

		maxConcurrency, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Incorrect maxConcurrency: %v %v\n", os.Args[2], err)
			os.Exit(1)
		}

		maxPages, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf("Incorrect maxPages: %v %v\n", os.Args[3], err)
			os.Exit(1)
		}
	}

	cfg := config{
		pages: make(map[string]int),
		baseURL: os.Args[1],
		mu: &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg: &sync.WaitGroup{},
		maxPages: maxPages,
	}

	fmt.Printf("starting crawl of: %v\n", os.Args[1])

	cfg.wg.Add(1)
	go cfg.crawlPage(os.Args[1])

	cfg.wg.Wait()

	fmt.Printf("================================\n")
	for k, v := range cfg.pages {
		fmt.Printf("%s: %d\n", k, v)
	}
}

