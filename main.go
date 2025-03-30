package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

const (
	ROOT = "https://scrape-me.dreamsofcode.io/"
)

func main() {
	results := make(map[string]int)
	scrape(ROOT, results)

	fmt.Println(results)

	// resp, err := http.Get("")
	// if err != nil {
	// 	log.Fatalf("Unable to open root website: %v", err)
	// }

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("Unable to read root body: %v", err)
	// }

	// sb := string(body)
	// log.Println(sb)
}

func scrape(link string, results map[string]int) {
	// Check if link already in
	if _, ok := results[link]; ok {
		return
	}

	r, err := http.Get(link)
	results[link] = r.StatusCode
	if err != nil {
		fmt.Println("Link is broken:", link)
		return
	}

	doc, err := html.Parse(r.Body)
	if err != nil {
		fmt.Println("Unable to parse body on link", link)
		return
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Println(attr.Val)
				}
			}
		}
	}
}
