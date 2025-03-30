package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

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
		fmt.Fprintf(w, "%v\t%v\n", link, status)
	}

	w.Flush()
}

func checkLink(link string, results map[string]int) {
	// Check if link already in
	if _, ok := results[link]; ok {
		return
	}

	fmt.Printf("Checking link %v...\n", link)

	r, err := http.Get(link)
	if err != nil {
		fmt.Println("Link is broken:", link)
		results[link] = -1
		return
	}

	defer r.Body.Close()

	// Store status code
	results[link] = r.StatusCode

	// Check status code
	if r.StatusCode > 400 {
		fmt.Printf("%v | Status: %v", link, r.Status)
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

	checkNode(doc)
}

func checkNode(node *html.Node) {
	for node != nil {
		if node.Type == html.ElementNode && node.Data == "a" {
			fmt.Println(node.Data)
		}

		// Recursively check each child
		checkNode(node.FirstChild)
		node = node.NextSibling
	}
}
