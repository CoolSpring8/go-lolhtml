# go-lolhtml

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/coolspring8/go-lolhtml/Go) ![Codecov](https://img.shields.io/codecov/c/github/coolspring8/go-lolhtml) [![Go Report Card](https://goreportcard.com/badge/github.com/coolspring8/go-lolhtml)](https://goreportcard.com/report/github.com/coolspring8/go-lolhtml) [![PkgGoDev](https://pkg.go.dev/badge/github.com/coolspring8/go-lolhtml)](https://pkg.go.dev/github.com/coolspring8/go-lolhtml)

Go bindings for the Rust library [cloudflare/lol-html](https://github.com/cloudflare/lol-html/), the *Low Output Latency streaming HTML rewriter/parser with CSS-selector based API*, talking via cgo.

**Status:** All abilities provided by C-API implemented, except for customized user data in handlers. Tests are partially covered. The code is at its early stage and the API is therefore subject to change. If you have any ideas on how API can be better structured, feel free to open a PR or an issue.

## Installation

For Linux/macOS/Windows x86_64 platforms, installation is as simple as a single `go get`:

```shell
$ go get github.com/coolspring8/go-lolhtml
```

There is no need for you to install Rust. That's because lol-html could be prebuilt into static libraries, stored and shipped in `/build` folder, so that cgo can handle other matters naturally and smoothly.

For other platforms, you'll have to compile it yourself.

## Getting Started

Now let's initialize a project and create `main.go`:

```go
package main

import (
	"bytes"
	"github.com/coolspring8/go-lolhtml"
	"io"
	"log"
	"os"
)

func main() {
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
	// Output: Hello, <span>LOL-HTML</span>!
}
```

The above program takes the chunk `Hello, <span>World</span>!` as input, is configured to rewrite all texts in `span` tags to "LOL-HTML" and prints the result to standard output.

And the result is `Hello, <span>LOL-HTML</span>!` .

For more examples, explore the `/examples` directory.

## Documentation

Available at [pkg.go.dev](https://pkg.go.dev/github.com/coolspring8/go-lolhtml). (WIP)

## Other Bindings

- Rust (native), C, JavaScript - [cloudflare/lol-html](https://github.com/cloudflare/lol-html/)
- Lua - [jdesgats/lua-lolhtml](https://github.com/jdesgats/lua-lolhtml/)

## License

BSD 3-Clause "New" or "Revised" License

## Disclaimer

This is an unofficial binding.

Cloudflare is a registered trademark of Cloudflare, Inc. Cloudflare names used in this project are for identification purposes only. The project is not associated in any way with Cloudflare Inc.