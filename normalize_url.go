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

	url := fmt.Sprintf("%s%s", parsedURL.Hostname(), urlPath)
	return url, nil

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
