package main

import (
//	"fmt"
	"strings"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/url"
)

func traverse(doc *html.Node) []*url.URL {
	links := []*url.URL{}
	if doc.Type == html.ElementNode && doc.DataAtom == atom.A {
		for _, a := range doc.Attr {
			if a.Key == "href" {
				newLink, err := url.Parse(a.Val)
				if err != nil {
					break
				}
				links = append(links, newLink)
				break
			}
		}
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, traverse(c)...)
	}		

	return links
}


func getURLSFromHTML(htmlBody string, bu *url.URL) ([]*url.URL, error) {

	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	links := traverse(doc)

	for i, l := range links {
		if l.Host == "" {
			// fmt.Printf("Pre: %v\n", l)
			links[i].Host = bu.Host
			links[i].Scheme = bu.Scheme
			// fmt.Printf("Post: %v\n", l)
		}
	}

	return links, nil
}

