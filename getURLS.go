package main

import (
//	"fmt"
	"strings"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func traverse(doc *html.Node) []string {
	links := []string{}
	if doc.Type == html.ElementNode && doc.DataAtom == atom.A {
		for _, a := range doc.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
				break
			}
		}
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, traverse(c)...)
	}		

	return links
}


func getURLSFromHTML(htmlBody, rawBaseURL string) ([]string, error) {

	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	links := traverse(doc)

	for i, l := range links {
		if strings.HasPrefix(l, "/") {
			links[i] = rawBaseURL + l
		}
	}

	return links, nil
}

