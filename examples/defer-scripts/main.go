// Usage: curl -NL https://git.io/JeOSZ | go run main.go
package main

import (
	"io"
	"log"
	"os"

	"github.com/coolspring8/go-lolhtml"
)

func main() {
	w, err := lolhtml.NewWriter(
		os.Stdout,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "script[src]:not([async]):not([defer])",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						err := e.SetAttribute("defer", "")
						if err != nil {
							log.Fatal(err)
						}
						return lolhtml.Continue
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Free()

	_, err = io.Copy(w, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	err = w.End()
	if err != nil {
		log.Fatal(err)
	}
}
