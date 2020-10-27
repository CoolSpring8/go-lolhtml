// This file is for demonstration in godoc. For more examples, see the /examples directory.
package lolhtml_test

import (
	"bytes"
	"github.com/coolspring8/go-lolhtml"
	"io"
	"log"
	"os"
)

func ExampleNewWriter() {
	chunk := []byte("Hello, <span>World</span>!")
	r := bytes.NewReader(chunk)
	w, err := lolhtml.NewWriter(
		os.Stdout,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "span",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						err := e.SetInnerContentAsRaw("LOL-HTML")
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

	_, err = io.Copy(w, r)
	if err != nil {
		log.Fatal(err)
	}

	err = w.End()
	if err != nil {
		log.Fatal(err)
	}
	// Output: Hello, <span>LOL-HTML</span>!
}
