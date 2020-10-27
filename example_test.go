// This file is for demonstration in godoc. For more examples, see the /examples directory.
package lolhtml_test

import (
	"bytes"
	"fmt"
	"github.com/coolspring8/go-lolhtml"
	"io"
	"log"
	"os"
	"strings"
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

func ExampleRewriteString() {
	output, err := lolhtml.RewriteString(
		`<div><a href="http://example.com"></a></div>`,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "a[href]",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						href, err := e.AttributeValue("href")
						if err != nil {
							log.Fatal(err)
						}
						href = strings.ReplaceAll(href, "http:", "https:")

						err = e.SetAttribute("href", href)
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

	fmt.Println(output)
	// Output: <div><a href="https://example.com"></a></div>
}
