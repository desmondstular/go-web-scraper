package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/net/html"
)

const (
	ROOT = "https://scrape-me.dreamsofcode.io"
)

func main() {
	results := make(map[string]int)
	checkLink(ROOT, results)

	// Tabular writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	fmt.Fprintf(w, "\nLink\tStatus\n")
	for link, status := range results {
		if status >= 400 {
			fmt.Fprintf(w, "%v\t%v\n", link, status)
		}
	}

	w.Flush()
}

func checkLink(link string, results map[string]int) {
	// Check if link already in
	if _, ok := results[link]; ok {
		return
	}

	client := &http.Client{
		Timeout: time.Second * 1,
		CheckRedirect: func(req *http.Request, voa []*http.Request) error {
			checkLink(req.URL.String(), results)
			return fmt.Errorf("Handling redirect")
		},
	}

	fmt.Printf("Checking link %v...\n", link)

	r, err := client.Get(link)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer r.Body.Close()

	// Store status code
	results[link] = r.StatusCode

	// Check status code
	if r.StatusCode >= 400 {
		return
	}

	// Do not parse link if it is a redirect; stay on the same site
	if !strings.HasPrefix(link, ROOT) {
		fmt.Println("Not parsing this link further:", link)
		return
	}

	doc, err := html.Parse(r.Body)
	if err != nil {
		fmt.Println("Unable to parse body on link", link)
		return
	}

	checkNode(doc, results)
}

func checkNode(node *html.Node, results map[string]int) {
	for node != nil {
		if node.Type == html.ElementNode && node.Data == "a" {
			// Get href attr
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					if strings.HasPrefix(attr.Val, "/") {
						checkLink(ROOT+attr.Val, results)
					} else {
						checkLink(attr.Val, results)
					}
				}
			}
		}

		// Check node's children recursively
		checkNode(node.FirstChild, results)

		node = node.NextSibling
	}
}
