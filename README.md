# Go Web Crawler

A simple concurrent web crawler written in Go.

The crawler starts from a given URL, follows links found on each page, extracts basic page information, and saves the results to a JSON report.

## Features

* Concurrent page crawling
* Configurable maximum concurrency
* Configurable maximum number of pages
* Restricts crawling to the starting URL's hostname
* Resolves relative and absolute links
* Normalizes URLs to avoid duplicate visits
* Extracts:

  * First `h1` or `h2` heading
  * First paragraph
  * Outgoing links
  * Image URLs
* Generates a formatted `report.json` file
* Uses `goquery` for HTML parsing

## Requirements

* Go 1.20 or later
* Internet connection

## Installation

Clone the repository:

```bash
git clone <repository-url>
cd webcrawler
```

Install dependencies:

```bash
go mod download
```

## Usage

Run the crawler with three arguments:

```bash
go run . <url> <max-concurrency> <max-pages>
```

Example:

```bash
go run . https://example.com 5 20
```

Arguments:

| Argument          | Description                                  |
| ----------------- | -------------------------------------------- |
| `url`             | Starting URL for the crawl                   |
| `max-concurrency` | Maximum number of pages crawled concurrently |
| `max-pages`       | Maximum number of pages to crawl             |

For example:

```bash
go run . https://example.com 10 100
```

This starts the crawler at `https://example.com`, allows up to 10 concurrent requests, and crawls a maximum of 100 pages.

## Output

The crawler creates a `report.json` file containing the extracted data for each crawled page.

Example:

```json
[
  {
    "URL": "https://example.com",
    "Heading": "Example Domain",
    "FirstParagraph": "This domain is for use in illustrative examples.",
    "OutgoingLinks": [
      "https://example.com/about"
    ],
    "ImageURLs": []
  }
]
```

Pages in the report are sorted by their normalized URL.

## How It Works

1. The starting URL is parsed and used as the base URL.
2. The crawler downloads the page using an HTTP GET request.
3. The HTML is parsed using `goquery`.
4. Links are extracted and resolved against the current page URL.
5. Page information is extracted.
6. New links are scheduled for crawling.
7. A concurrency channel limits the number of active crawlers.
8. Previously visited URLs are skipped.
9. Crawling stops when the maximum page limit is reached.
10. The collected data is written to `report.json`.

## Project Structure

```text
.
├── main.go
├── crawler.go
├── html.go
├── report.go
├── config.go
├── go.mod
└── report.json
```

The exact file names may differ depending on how the project is organized.

## Dependencies

This project uses:

* `github.com/PuerkitoBio/goquery` for HTML parsing

## Learning Goals

This project was built to practice Go concepts including:

* Goroutines
* Channels
* `sync.WaitGroup`
* `sync.Mutex`
* HTTP requests
* URL parsing and resolution
* HTML parsing
* Concurrent crawling
* JSON serialization
* Command-line arguments
* File I/O
* Basic concurrency control

