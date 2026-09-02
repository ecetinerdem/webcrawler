package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func normalizeURL(urlString string) (string, error) {

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	urlPath := parsedURL.Path

	urlPath = strings.TrimSuffix(urlPath, "/")

	resultURL := fmt.Sprintf("%s%s", parsedURL.Hostname(), urlPath)
	return resultURL, nil

}

func getHeadingFromHTML(html string) string {
	reader := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Println("Error goquery: ", err)
		return ""
	}

	hOne := doc.Find("h1")
	if hOne.Length() > 0 {
		return hOne.Eq(0).Text()
	}
	hTwo := doc.Find("h2")
	if hTwo.Length() > 0 {
		return hTwo.Eq(0).Text()
	}

	log.Println("Could not find h1 or h2")
	return ""
}

func getFirstParagraphFromHTML(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Println("Error goquery: ", err)
		return ""
	}

	paragraph := doc.Find("main").First().Find("p").First()
	if paragraph.Length() == 0 {
		paragraph = doc.Find("p")
		if paragraph.Length() == 0 {
			return ""
		}
		return paragraph.Eq(0).Text()
	}

	return paragraph.Eq(0).Text()
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	var urls []string

	if err != nil {
		log.Println("Error goquery: ", err)
		return urls, err
	}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		relativePath, ok := s.Attr("href")
		if !ok {
			log.Println("Err in relative path, path may not exist")
			return
		}

		urlFromRelative, err := url.Parse(relativePath)
		if err != nil {
			log.Println("Err in relative path: ", err)
			return
		}

		resolved := baseURL.ResolveReference(urlFromRelative)

		urls = append(urls, resolved.String())

	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	var urls []string

	if err != nil {
		log.Println("Error goquery: ", err)
		return urls, err
	}

	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		relativePath, ok := s.Attr("src")
		if !ok {
			log.Println("Err in relative path: ", err)
			return
		}

		urlFromRelative, err := url.Parse(relativePath)
		if err != nil {
			log.Println("Err in relative path: ", err)
			return
		}

		resolved := baseURL.ResolveReference(urlFromRelative)

		urls = append(urls, resolved.String())

	})

	return urls, nil

}

func extractPageData(html string, pageURL string) PageData {

	var pageData PageData

	heading := getHeadingFromHTML(html)
	firstParagraph := getFirstParagraphFromHTML(html)
	pageData.Heading = heading
	pageData.FirstParagraph = firstParagraph

	urlStruct, err := url.Parse(pageURL)
	if err != nil {
		return pageData
	}
	pageData.URL = urlStruct.String()

	outgoingLinks, err := getURLsFromHTML(html, urlStruct)
	if err != nil {
		return pageData
	}
	pageData.OutgoingLinks = outgoingLinks

	imageURLs, err := getImagesFromHTML(html, urlStruct)
	if err != nil {
		return pageData
	}
	pageData.ImageURLs = imageURLs

	return pageData
}
