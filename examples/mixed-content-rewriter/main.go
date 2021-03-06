// Usage: curl -NL https://git.io/JeOSZ | go run main.go
package main

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/coolspring8/go-lolhtml"
)

func main() {
	w, err := lolhtml.NewWriter(
		os.Stdout,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "a[href], link[rel=stylesheet][href]",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						rewriteUrlInAttribute(e, "href")
						return lolhtml.Continue
					},
				},
				{
					Selector: "script[src], iframe[src], img[src], audio[src], video[src]",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						rewriteUrlInAttribute(e, "src")
						return lolhtml.Continue
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(w, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func rewriteUrlInAttribute(e *lolhtml.Element, attributeName string) {
	attr, err := e.AttributeValue(attributeName)
	if err != nil {
		log.Fatal(err)
	}
	attr = strings.ReplaceAll(attr, "http://", "https://")

	err = e.SetAttribute(attributeName, attr)
	if err != nil {
		log.Fatal(err)
	}
}
