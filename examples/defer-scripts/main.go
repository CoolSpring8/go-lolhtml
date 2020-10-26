package main

import (
	"github.com/coolspring8/go-lolhtml"
	"io"
	"log"
	"os"
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
}
