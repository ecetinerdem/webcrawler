package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getHTML(rawURL string) (string, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", rawURL, nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("status code is above 400")
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("content type is not text/html")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	html := string(body)
	return html, nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {

	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if cfg.maxPages <= len(cfg.pages) {
		return
	}
	rawCurrent, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Println("error parsing current: ", err)
		return
	}

	if cfg.baseURL.Hostname() != rawCurrent.Hostname() {
		return
	}

	normalizedRawURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	isFirst := cfg.addPageVisit(normalizedRawURL)

	if !isFirst {
		return
	}

	fmt.Printf("crawling %s\n", rawCurrentURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	urls, err := getURLsFromHTML(html, rawCurrent)
	if err != nil {
		return
	}

	pageData := extractPageData(html, rawCurrentURL)

	cfg.pages[normalizedRawURL] = pageData

	for _, u := range urls {
		cfg.wg.Add(1)
		go func() {
			newU := u
			cfg.crawlPage(newU)
		}()
	}

}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	isFirst = true

	_, exist := cfg.pages[normalizedURL]

	if exist {
		return !isFirst
	}

	cfg.pages[normalizedURL] = PageData{}
	return isFirst
}
