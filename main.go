package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {

	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	rawURL := args[0]
	fmt.Printf("starting crawl: %s\n", rawURL)

	maxConcurency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("max concurrency must be integer")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("max pages must be integer")
		os.Exit(1)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("cannot parse url")
		os.Exit(1)
	}
	cfg := config{}
	cfg.pages = make(map[string]PageData)
	cfg.baseURL = parsedURL
	cfg.mu = &sync.Mutex{}
	cfg.wg = &sync.WaitGroup{}
	cfg.concurrencyControl = make(chan struct{}, maxConcurency)
	cfg.maxPages = maxPages

	cfg.wg.Add(1)
	go cfg.crawlPage(cfg.baseURL.String())

	cfg.wg.Wait()

	for k, v := range cfg.pages {
		fmt.Printf("%s: %v\n", k, v)
	}

	err = writeJSONReport(cfg.pages, "report.json")
	if err != nil {
		fmt.Println("could not write report: ", err)
		os.Exit(1)
	}

}
